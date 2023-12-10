package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleAuthorDao interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

type ArticleGORMAuthorDAO struct {
	db *gorm.DB
}

func (a ArticleGORMAuthorDAO) Create(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a ArticleGORMAuthorDAO) Update(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

func NewArticleGORMAuthorDAO(db *gorm.DB) ArticleAuthorDao {
	return &ArticleGORMAuthorDAO{
		db: db,
	}
}
