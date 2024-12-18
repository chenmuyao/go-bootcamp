package service

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
)

var ErrArticleNotFound = repository.ErrArticleNotFound

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.ID > 0 {
		return article.ID, a.repo.Update(ctx, article)
	} else {
		return a.repo.Create(ctx, article)
	}
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	if article.ID > 0 {
		return article.ID, a.repo.Update(ctx, article)
	} else {
		return a.repo.Create(ctx, article)
	}
}
