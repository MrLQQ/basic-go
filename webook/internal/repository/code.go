package repository

import (
	"basic-go/webook/internal/repository/cache"
	"context"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inCode string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCacheCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (c *CacheCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CacheCodeRepository) Verify(ctx context.Context, biz, phone, inCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inCode)
}
