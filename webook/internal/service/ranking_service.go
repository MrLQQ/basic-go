package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"log"
	"time"
)

type RankingService interface {
	// TopN 前100
	TopN(ctx context.Context) error
}

type BatchRankingService struct {
	// 用来取点赞数
	intrSvc InteractiveService

	// 用来查找文章
	artSvc ArticleService

	batchSize int
	scoreFunc func(likeCnt int64, utime time.Time) float64
	n         int
}

func NewBatchRankingService(intrSvc InteractiveService, artSvc ArticleService) *BatchRankingService {
	return &BatchRankingService{intrSvc: intrSvc, artSvc: artSvc}
}

func (b *BatchRankingService) TopN(ctx context.Context) error {
	arts, err := b.topN(ctx)
	if err != nil {
		return nil
	}
	// 最终是要放到缓存里面的
	// 存到缓存里面
	log.Println(arts)
	panic("implement me")
}

func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")

}
