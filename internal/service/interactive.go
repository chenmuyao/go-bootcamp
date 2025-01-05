package service

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	CancelCollect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
}

type interactiveService struct {
	repo repository.InteractiveRepository
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
		repo: repo,
	}
}
