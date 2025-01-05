package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func (c *CachedInteractiveRepository) DecrLike(
	ctx context.Context,
	biz string,
	id int64,
	uid int64,
) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}

	return c.cache.DecrLikeCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) IncrLike(
	ctx context.Context,
	biz string,
	id int64,
	uid int64,
) error {
	err := c.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}

	return c.cache.IncrLikeCntIfPresent(ctx, biz, id)
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
