package cache

import (
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
}

type RankingRedisCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func NewRankingRedisCache(client redis.Cmdable) RankingCache {
	return &RankingRedisCache{client: client,
		key:        "ranking:top_n",
		expiration: time.Minute * 3,
	}
}

func (r *RankingRedisCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}
