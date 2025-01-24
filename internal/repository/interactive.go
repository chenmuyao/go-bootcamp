package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	DeleteCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizID int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, bizID int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, bizID int64, uid int64) (bool, error)
}

type CachedInteractiveRepository struct {
	l     logger.Logger
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

// Collected implements InteractiveRepository.
func (c *CachedInteractiveRepository) Collected(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, bizID, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

// Get implements InteractiveRepository.
func (c *CachedInteractiveRepository) Get(
	ctx context.Context,
	biz string,
	bizID int64,
) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizID)
	if err == nil {
		slog.Error("intr cache", slog.Any("intr", intr))
		return intr, nil
	}
	intrDAO, err := c.dao.Get(ctx, biz, bizID)
	if err != nil {
		return domain.Interactive{}, nil
	}
	slog.Error("intr dao", slog.Any("intrDAO", intrDAO))
	res := c.toDomain(intrDAO)
	err = c.cache.Set(ctx, biz, bizID, res)
	if err != nil {
		c.l.Error(
			"failed to set interactive cache",
			logger.String("biz", biz),
			logger.Int64("bizID", bizID),
			logger.Error(err),
		)
	}
	slog.Error("intr res", slog.Any("intr", res))
	return res, nil
}

// Liked implements InteractiveRepository.
func (c *CachedInteractiveRepository) Liked(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, bizID, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

// DeleteCollectionItem implements InteractiveRepository.
func (c *CachedInteractiveRepository) DeleteCollectionItem(
	ctx context.Context,
	biz string,
	id int64,
	cid int64,
	uid int64,
) error {
	err := c.dao.DeleteCollectionBiz(ctx, dao.UserCollectionBiz{
		UID:   uid,
		BizID: id,
		Biz:   biz,
	})
	if err != nil {
		return err
	}

	return c.cache.DecrCollectCntIfPresent(ctx, biz, id)
}

// AddCollectionItem implements InteractiveRepository.
func (c *CachedInteractiveRepository) AddCollectionItem(
	ctx context.Context,
	biz string,
	id int64,
	cid int64,
	uid int64,
) error {
	now := time.Now().UnixMilli()
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		UID:   uid,
		BizID: id,
		Biz:   biz,
		CID:   cid,
		Utime: now,
		Ctime: now,
	})
	if err != nil {
		return err
	}

	return c.cache.IncrCollectCntIfPresent(ctx, biz, id)
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

func (c *CachedInteractiveRepository) toDomain(dao dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    dao.ReadCnt,
		LikeCnt:    dao.LikeCnt,
		CollectCnt: dao.CollectCnt,
	}
}

func NewCachedInteractiveRepository(
	l logger.Logger,
	dao dao.InteractiveDAO,
	cache cache.InteractiveCache,
) InteractiveRepository {
	return &CachedInteractiveRepository{
		l:     l,
		dao:   dao,
		cache: cache,
	}
}
