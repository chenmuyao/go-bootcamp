package service

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
}

type interactiveService struct {
	repo repository.InteractiveRepository
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
