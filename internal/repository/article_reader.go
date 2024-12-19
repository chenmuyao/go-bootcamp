package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

type ArticleReaderRepository interface {
	// Insert or Update
	Save(ctx context.Context, article domain.Article) error
}

type CachedArticleReaderRepository struct {
	dao dao.ArticleDAO
}

func NewArticleReaderRepository(dao dao.ArticleDAO) ArticleReaderRepository {
	return &CachedArticleReaderRepository{
		dao: dao,
	}
}

func (c *CachedArticleReaderRepository) Save(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateByID(ctx, c.toEntity(article))
}

// func (c *CachedArticleReaderRepository) toDAO(article dao.Article) domain.Article {
// 	return domain.Article{
// 		Title:   article.Title,
// 		Content: article.Content,
// 	}
// }

func (c *CachedArticleReaderRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		ID:       article.ID,
		Title:    article.Title,
		Content:  article.Content,
		AuthorID: article.Author.ID,
	}
}
