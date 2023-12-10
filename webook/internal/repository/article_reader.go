package repository

import (
	"basic-go/webook/internal/domain"
	"context"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}
