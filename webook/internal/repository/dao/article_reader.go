package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDao interface {
	// Upsert INSERT or UPDATE
	Upsert(ctx context.Context, art Article) error
	UpsertV2(ctx context.Context, art PublishArticle) error
}

type ArticleGORMReaderDAO struct {
	db *gorm.DB
}

func (a ArticleGORMReaderDAO) UpsertV2(ctx context.Context, art PublishArticle) error {
	//TODO implement me
	panic("implement me")
}

func NewArticleGORMReaderDAO(db *gorm.DB) ArticleReaderDao {
	return &ArticleGORMReaderDAO{db: db}
}

func (a ArticleGORMReaderDAO) Upsert(ctx context.Context, art Article) error {

	//TODO implement me
	panic("implement me")
}
