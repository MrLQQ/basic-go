package repository

import (
	"basic-go/webook/internal/repository/cache"
	"context"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

// CodeRedisRepository ---------------------------------------使用Redis缓存实现--------------------------------------------------//
type CodeRedisRepository struct {
	cache *cache.CodeRedisCache
}

func NewCodeRedisRepository(c *cache.CodeRedisCache) *CodeRedisRepository {
	return &CodeRedisRepository{
		cache: c,
	}
}

func (c *CodeRedisRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRedisRepository) Verify(ctx context.Context, biz, phone, inCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inCode)
}

// CodeCacheRepository ---------------------------------------使用本地缓存实现--------------------------------------------------//
type CodeCacheRepository struct {
	cache *cache.CodeCache
}

func NewCodeCacheRepository(c *cache.CodeCache) *CodeCacheRepository {
	return &CodeCacheRepository{
		cache: c,
	}
}

func (c *CodeCacheRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeCacheRepository) Verify(ctx context.Context, biz, phone, inCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inCode)
}
