package dao

import (
	"context"

	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	// Insert and Update
	Upsert(ctx context.Context, article Article) error
	UpsertV2(ctx context.Context, article PublishedArticle) error
}

type ArticleGORMReaderDAO struct {
	db *gorm.DB
}

func (a *ArticleGORMReaderDAO) UpsertV2(ctx context.Context, article PublishedArticle) error {
	panic("unimplemented")
}

func (a *ArticleGORMReaderDAO) Upsert(ctx context.Context, article Article) error {
	panic("unimplemented")
}

func NewArticleGORMReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &ArticleGORMReaderDAO{
		db: db,
	}
}
