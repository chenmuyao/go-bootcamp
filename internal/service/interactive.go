package service

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -source=./interactive.go -package=svcmocks -destination=./mocks/interactive.mock.go
type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	CancelCollect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
	// NOTE: Intr must exist
	MustBatchGet(ctx context.Context, biz string, id []int64) ([]domain.Interactive, error)
	GetByIDs(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error)
	GetTopLike(ctx context.Context, biz string, limit int) ([]int64, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository

	defaultTopLikeLimit int
}

// GetByIDs implements InteractiveService.
func (i *interactiveService) GetByIDs(
	ctx context.Context,
	biz string,
	ids []int64,
) (map[int64]domain.Interactive, error) {
	intrs, err := i.repo.GetByIDs(ctx, biz, ids)
	if err != nil {
		return nil, err
	}

	res := make(map[int64]domain.Interactive)
	for _, intr := range intrs {
		res[intr.BizID] = intr
	}
	return res, nil
}

// GetTopLike implements InteractiveService.
func (i *interactiveService) GetTopLike(
	ctx context.Context,
	biz string,
	limit int,
) ([]int64, error) {
	if limit <= 0 || limit > i.defaultTopLikeLimit {
		limit = i.defaultTopLikeLimit
	}
	return i.repo.GetTopLike(ctx, biz, limit)
}

// BatchGet implements InteractiveService.
func (i *interactiveService) MustBatchGet(
	ctx context.Context,
	biz string,
	id []int64,
) ([]domain.Interactive, error) {
	return i.repo.MustBatchGet(ctx, biz, id)
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
