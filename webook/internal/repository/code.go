package repository

import (
	"basic-go/webook/internal/repository/cache"
	"context"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

type CodeRepository struct {
	cache cache.CodeCache
}

func (c *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) Verify(ctx context.Context, biz, phone, inCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inCode)
}
