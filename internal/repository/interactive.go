package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

// IncrReadCnt implements InteractiveRepository.
func (c *CachedInteractiveRepository) IncrReadCnt(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizID)
	if err != nil {
		return err
	}

	// NOTE: add cache might fail and cause the inconsistency of data, but it
	// is not critical in this feature.
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizID)
}

func NewCachedInteractiveRepository(
	dao dao.InteractiveDAO,
	cache cache.InteractiveCache,
) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		cache: cache,
	}
}
