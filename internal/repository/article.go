package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDAO
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) Create(
	ctx context.Context,
	article domain.Article,
) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(article))
}

// func (c *CachedArticleRepository) toDAO(article dao.Article) domain.Article {
// 	return domain.Article{
// 		Title:   article.Title,
// 		Content: article.Content,
// 	}
// }

func (c *CachedArticleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorID: article.Author.ID,
	}
}
