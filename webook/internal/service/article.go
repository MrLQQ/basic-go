package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"basic-go/webook/pkg/logger"
	"context"
	"errors"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, id int64) error
}

type articleService struct {
	repo repository.ArticleRepository

	//	V1写法专用
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
	l          logger.LoggerV1
}

func (a *articleService) Withdraw(ctx context.Context, uid int64, id int64) error {
	return a.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {

	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		// 我可能线上库已经有数据了
		// 也可能没有
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			a.l.Error("保存到制作库成功，但是到线上库失败",
				logger.Int64("aid", art.Id),
				logger.Error(err))
		} else {
			return id, nil
		}

	}
	a.l.Error("保存到制作库成功，但是到线上库失败，重试次数耗尽",
		logger.Int64("aid", art.Id),
		logger.Error(err))
	return id, errors.New("保存到线上库失败，重试次数耗尽")

}

func NewArticleServiceV1(
	readerRepo repository.ArticleReaderRepository,
	authorRepo repository.ArticleAuthorRepository,
	l logger.LoggerV1) *articleService {
	return &articleService{
		readerRepo: readerRepo,
		authorRepo: authorRepo,
		l:          l,
	}
}
func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	} else {
		return a.repo.Create(ctx, art)
	}

}
