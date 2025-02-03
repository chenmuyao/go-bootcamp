package service

import (
	"context"
	"errors"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/events/article"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

const publishMaxRetry = 3

var (
	ErrArticleNotFound = repository.ErrArticleNotFound
	ErrPublish         = errors.New("still failed to publish article after retries")
)

//go:generate mockgen -source=./article.go -package=svcmocks -destination=./mocks/article.mock.go
type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, userID int64, articleID int64) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPubByID(ctx context.Context, id int64, uid int64) (domain.Article, error)
	BatchGetPubByIDs(ctx context.Context, ids []int64) ([]domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error)
}

type articleService struct {
	l    logger.Logger
	repo repository.ArticleRepository

	// v1: separate reader and author at repo level
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
	producer   article.Producer
}

// ListPub implements ArticleService.
func (a *articleService) ListPub(
	ctx context.Context,
	start time.Time,
	offset int,
	limit int,
) ([]domain.Article, error) {
	panic("unimplemented")
}

// BatchGetPubByIDs implements ArticleService.
func (a *articleService) BatchGetPubByIDs(
	ctx context.Context,
	ids []int64,
) ([]domain.Article, error) {
	return a.repo.BatchGetPubByIDs(ctx, ids)
}

// GetPubByID implements ArticleService.
func (a *articleService) GetPubByID(
	ctx context.Context,
	id int64,
	uid int64,
) (domain.Article, error) {
	res, err := a.repo.GetPubByID(ctx, id)
	go func() {
		if err == nil {
			// send a message
			er := a.producer.ProduceReadEvent(article.ReadEvent{
				Aid: id,
				Uid: uid,
			})
			if er != nil {
				a.l.Error(
					"failed to send ReadEvent",
					logger.Int64("aid", id),
					logger.Int64("uid", uid),
					logger.Error(err),
				)
			}
		}
	}()
	return res, err
}

// GetByID implements ArticleService.
func (a *articleService) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetByID(ctx, id)
}

// GetByAuthor implements ArticleService.
func (a *articleService) GetByAuthor(
	ctx context.Context,
	uid int64,
	offset int,
	limit int,
) ([]domain.Article, error) {
	return a.repo.GetByAuthor(ctx, uid, offset, limit)
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnpublished
	if article.ID > 0 {
		return article.ID, a.repo.Update(ctx, article)
	} else {
		return a.repo.Create(ctx, article)
	}
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, article)
}

func (a *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished
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

func (a *articleService) Withdraw(ctx context.Context, userID int64, articleID int64) error {
	return a.repo.SyncStatus(ctx, userID, articleID, domain.ArticleStatusPrivate)
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

func NewArticleService(
	repo repository.ArticleRepository,
	producer article.Producer,
) ArticleService {
	return &articleService{
		repo:     repo,
		producer: producer,
	}
}
