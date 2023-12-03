package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, entity Article) error
}

type ArticleGORMDAO struct {
	db *gorm.DB
}

func (a *ArticleGORMDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).Model(&art).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败，ID不会或者作者不对")
	}
	return nil
}

func NewArticleGORMDAO(db *gorm.DB) ArticleDAO {
	return &ArticleGORMDAO{
		db: db,
	}
}

func (a *ArticleGORMDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := a.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type=BLOB"`
	// 我要根据创作者ID来查询
	AuthorId int64 `gorm:"index"`
	Status   uint8 `bson:"status,omitempty"`
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
