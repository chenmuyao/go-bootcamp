package service

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

const publishMaxRetry = 3

var (
	ErrArticleNotFound = repository.ErrArticleNotFound
	ErrPublish         = errors.New("still failed to publish article after retries")
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	l    logger.Logger
	repo repository.ArticleRepository

	// v1: separate reader and author at repo level
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
}

func NewArticleServiceV1(
	l logger.Logger,
	readerRepo repository.ArticleReaderRepository,
	authorRepo repository.ArticleAuthorRepository,
) *articleService {
	return &articleService{
		l:          l,
		readerRepo: readerRepo,
		authorRepo: authorRepo,
	}
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
	return a.repo.Sync(ctx, article)
}

func (a *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	// first change author repo
	id := article.ID
	var err error

	if id > 0 {
		err = a.authorRepo.Update(ctx, article)
	} else {
		id, err = a.authorRepo.Create(ctx, article)
		if err != nil {
			return 0, err
		}
	}

	// then change reader repo
	article.ID = id
	for range publishMaxRetry {
		err = a.readerRepo.Save(ctx, article)
		if err != nil {
			// NOTE: At service level, we don't know what are the destination of repositories,
			// it can be SQL, NoSQL, or S3 etc.
			// So we are not supposed to open a "transaction". And if we use 2 DB, only if we
			// use distributed DB, otherwise we cannot open a transaction at all.
			a.l.Error("Articled saved to author repo, but failed to publish to reader repo",
				logger.Int64("aid", id),
				logger.Error(err))
		} else {
			return id, nil
		}
	}
	return id, ErrPublish
}
