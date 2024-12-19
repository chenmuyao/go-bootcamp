package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrArticleNotFound = errors.New("article not found")

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateByID(ctx context.Context, article Article) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

type Article struct {
	ID      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(4096)"`
	Content string `gorm:"type=BLOB"`

	AuthorID int64 `gorm:"index"`
	Status   uint8
	Ctime    int64
	Utime    int64
}

// same DB, different tables
type PublishedArticle Article

func (a *GORMArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := a.db.WithContext(ctx).Create(&article).Error
	return article.ID, err
}

func (a *GORMArticleDAO) UpdateByID(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).
		Model(&Article{}).
		Where("id = ? AND author_id = ?", article.ID, article.AuthorID).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrArticleNotFound
	}
	return nil
}
