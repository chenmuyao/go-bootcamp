package dao

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockgen -source=./article_author.go -package=daomocks -destination=./mocks/article_author.mock.go
type ArticleAuthorDAO interface {
	Create(ctx context.Context, article Article) (int64, error)
	UpdateByID(ctx context.Context, article Article) error
}

type ArticleGORMAuthorDAO struct {
	db *gorm.DB
}

func NewArticleGORMAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	return &ArticleGORMAuthorDAO{
		db: db,
	}
}

func (a *ArticleGORMAuthorDAO) Create(ctx context.Context, article Article) (int64, error) {
	panic("unimplemented")
}

func (a *ArticleGORMAuthorDAO) UpdateByID(ctx context.Context, article Article) error {
	panic("unimplemented")
}
