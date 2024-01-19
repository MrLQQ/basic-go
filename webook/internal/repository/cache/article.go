package cache

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/dao"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, res []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	Set(ctx context.Context, art dao.Article) error
}

type ArticleRedisCache struct {
	client redis.Cmdable
}

func (a ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	key := a.firstKey(uid)
	//val, err := a.client.Get(ctx, firstKey).Result()
	val, err := a.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(val, &res)
	return nil, err
}

func (a ArticleRedisCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	// 不能把所有的都缓存
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}
	key := a.firstKey(uid)
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, key, val, time.Minute*10).Err()
}

func (*ArticleRedisCache) firstKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
