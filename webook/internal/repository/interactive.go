package repository

import (
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, id)
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, cache: cache}
}

func (c *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	// 要更新缓存了
	// 部分失败问题 --  数据不一致
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}
