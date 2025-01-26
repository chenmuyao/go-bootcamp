package service

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	CancelCollect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	GetTopLike(ctx context.Context, biz string, limit int) ([]domain.ArticleInteractive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository

	defaultTopLikeLimit int
}

// GetTopLike implements InteractiveService.
func (i *interactiveService) GetTopLike(
	ctx context.Context,
	biz string,
	limit int,
) ([]domain.ArticleInteractive, error) {
	if limit <= 0 || limit > i.defaultTopLikeLimit {
		limit = i.defaultTopLikeLimit
	}
	return i.repo.GetTopLike(ctx, biz, limit)
}

// Get implements InteractiveService.
func (i *interactiveService) Get(
	ctx context.Context,
	biz string,
	id int64,
	uid int64,
) (domain.Interactive, error) {
	intr, err := i.repo.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	// NOTE: can consider degrading
	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		intr.Liked, er = i.repo.Liked(ctx, biz, id, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		intr.Collected, er = i.repo.Collected(ctx, biz, id, uid)
		return er
	})
	return intr, eg.Wait()
}

// CancelCollect implements InteractiveService.
func (i *interactiveService) CancelCollect(
	ctx context.Context,
	biz string,
	id int64,
	cid int64,
	uid int64,
) error {
	return i.repo.DeleteCollectionItem(ctx, biz, id, cid, uid)
}

// Collect implements InteractiveService.
func (i *interactiveService) Collect(
	ctx context.Context,
	biz string,
	id int64,
	cid int64,
	uid int64,
) error {
	return i.repo.AddCollectionItem(ctx, biz, id, cid, uid)
}

// CancelLike implements InteractiveService.
func (i *interactiveService) CancelLike(
	ctx context.Context,
	biz string,
	id int64,
	uid int64,
) error {
	return i.repo.DecrLike(ctx, biz, id, uid)
}

// Like implements InteractiveService.
func (i *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.IncrLike(ctx, biz, id, uid)
}

// IncrReadCnt implements InteractiveService.
func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	return i.repo.IncrReadCnt(ctx, biz, bizID)
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo:                repo,
		defaultTopLikeLimit: 10,
	}
}
