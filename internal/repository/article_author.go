package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

//go:generate mockgen -source=./article_author.go -package=repomocks -destination=./mocks/article_author.mock.go
type ArticleAuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}

type CachedArticleAuthorRepository struct {
	dao dao.ArticleDAO
}

func NewArticleAuthorRepository(dao dao.ArticleDAO) ArticleAuthorRepository {
	return &CachedArticleAuthorRepository{
		dao: dao,
	}
}

func (c *CachedArticleAuthorRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateByID(ctx, c.toEntity(article))
}

func (c *CachedArticleAuthorRepository) Create(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(article))
}

// func (c *CachedArticleAuthorRepository) toDAO(article dao.Article) domain.Article {
// 	return domain.Article{
// 		Title:   article.Title,
// 		Content: article.Content,
// 	}
// }

func (c *CachedArticleAuthorRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		ID:       article.ID,
		Title:    article.Title,
		Content:  article.Content,
		AuthorID: article.Author.ID,
	}
}
