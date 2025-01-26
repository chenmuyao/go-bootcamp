package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/chenmuyao/generique/gslice"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIDs []int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	DeleteCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	BatchGet(ctx context.Context, biz string, bizIDs []int64) ([]domain.Interactive, error)
	Get(ctx context.Context, biz string, bizID int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, bizID int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, bizID int64, uid int64) (bool, error)
	GetTopLike(ctx context.Context, biz string, limit int) ([]domain.ArticleInteractive, error)
}

type CachedInteractiveRepository struct {
	l                    logger.Logger
	dao                  dao.InteractiveDAO
	cache                cache.InteractiveCache
	topCache             cache.TopArticlesCache
	articleRepo          ArticleRepository
	defaultTopLikedLimit int64
}

// GetTopLike implements InteractiveRepository.
func (c *CachedInteractiveRepository) GetTopLike(
	ctx context.Context,
	biz string,
	limit int,
) ([]domain.ArticleInteractive, error) {
	// Get top like articles' IDs from local cache
	res, err := c.topCache.GetTopLikedArticles(ctx)
	if err == nil && len(res) > 0 {
		if len(res) > limit {
			return res[:limit], nil
		}
		return res, nil
	}

	// If not found, compute from redis
	ids, err := c.cache.GetTopLikedIDs(ctx, biz, c.defaultTopLikedLimit)
	if err != nil {
		// XXX: The data should be prepared, if not found,
		// just return error
		c.l.Error("failed to get top liked ids", logger.String("biz", biz), logger.Error(err))
		return []domain.ArticleInteractive{}, err
	}

	// get articles and intr to construct the result

	articles, err := c.articleRepo.BatchGetPubByIDs(ctx, ids)
	if err != nil {
		c.l.Error("failed to get articles", logger.String("biz", biz), logger.Error(err))
		return []domain.ArticleInteractive{}, err
	}

	intrs, err := c.BatchGet(ctx, biz, ids)
	if err != nil {
		c.l.Error("failed to get interactions", logger.String("biz", biz), logger.Error(err))
		return []domain.ArticleInteractive{}, err
	}

	artInts := make([]domain.ArticleInteractive, 0, len(ids))
	for i := range ids {
		artInt := domain.ArticleInteractive{
			Article: articles[i],
			Intr:    intrs[i],
		}
		artInts = append(artInts, artInt)
	}

	// set back to local cache with a short ttl
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := c.topCache.SetTopLikedArticles(ctx, artInts)
		if er != nil {
			c.l.Error("failed to write back to local cache", logger.Error(err))
		}
	}()

	return artInts, nil
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

// BatchGet implements InteractiveRepository.
func (c *CachedInteractiveRepository) BatchGet(
	ctx context.Context,
	biz string,
	bizIDs []int64,
) ([]domain.Interactive, error) {
	intrs, err := c.cache.BatchGet(ctx, biz, bizIDs)
	if err == nil {
		slog.Error("intr cache", slog.Any("intr", intrs))
		return intrs, nil
	}
	intrDAOs, err := c.dao.BatchGet(ctx, biz, bizIDs)
	if err != nil {
		return []domain.Interactive{}, nil
	}
	res := gslice.Map(intrDAOs, func(id int, src dao.Interactive) domain.Interactive {
		return c.toDomain(src)
	})
	err = c.cache.BatchSet(ctx, biz, bizIDs, res)
	if err != nil {
		c.l.Error(
			"failed to set interactive cache",
			logger.String("biz", biz),
			logger.Any("bizID", bizIDs),
			logger.Error(err),
		)
	}
	return res, nil
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

// BatchIncrReadCnt implements InteractiveRepository.
func (c *CachedInteractiveRepository) BatchIncrReadCnt(
	ctx context.Context,
	bizs []string,
	bizIDs []int64,
) error {
	err := c.dao.BatchIncrReadCnt(ctx, bizs, bizIDs)
	if err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		for i, bizID := range bizIDs {
			er := c.cache.IncrReadCntIfPresent(ctx, bizs[i], bizID)
			if er != nil {
				c.l.Error(
					"failed to incr ReadCnt cache",
					logger.String("biz", bizs[i]),
					logger.Int64("bizID", bizID),
					logger.Error(er),
				)
			}
		}
	}()
	return nil
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
	topCache cache.TopArticlesCache,
	articleRepo ArticleRepository,
) InteractiveRepository {
	return &CachedInteractiveRepository{
		l:                    l,
		dao:                  dao,
		cache:                cache,
		topCache:             topCache,
		articleRepo:          articleRepo,
		defaultTopLikedLimit: 10,
	}
}
