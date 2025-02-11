package rediscache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "embed"

	"github.com/chenmuyao/generique/gslice"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

const (
	fieldReadCnt    = "read_cnt"
	fieldLikeCnt    = "like_cnt"
	fieldCollectCnt = "collect_cnt"
	intrExpiryTime  = time.Minute * 15
)

//go:embed lua/incr_cnt.lua
var luaIncrCnt string

//go:embed lua/incr_rank.lua
var luaIncrRank string

type InteractiveRedisCache struct {
	cache.BaseInteractiveCache
	client redis.Cmdable
}

// IncrLikeRank implements cache.InteractiveCache.
func (i *InteractiveRedisCache) IncrLikeRank(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	return i.client.ZIncrBy(ctx, i.topLikedKey(biz), float64(1), strconv.FormatInt(bizID, 10)).Err()
}

// DecrLikeRank implements cache.InteractiveCache.
func (i *InteractiveRedisCache) DecrLikeRank(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	return i.client.ZIncrBy(ctx, i.topLikedKey(biz), float64(-1), strconv.FormatInt(bizID, 10)).
		Err()
}

// GetTopLikedIDs implements cache.InteractiveCache.
func (i *InteractiveRedisCache) GetTopLikedIDs(
	ctx context.Context,
	biz string,
	limit int64,
) ([]int64, error) {
	resStr, err := i.client.ZRevRange(ctx, i.topLikedKey(biz), 0, limit).Result()
	if err != nil {
		return []int64{}, err
	}
	res := gslice.Map(resStr, func(id int, src string) int64 {
		i, _ := strconv.ParseInt(src, 10, 64)
		// parsing error ignored
		return i
	})
	return res, nil
}

// SetLikeToZSET implements cache.InteractiveCache.
func (i *InteractiveRedisCache) SetLikeToZSET(
	ctx context.Context,
	biz string,
	bizId int64,
	likeCnt int64,
) error {
	return i.client.ZAdd(ctx, i.topLikedKey(biz), redis.Z{
		Score:  float64(likeCnt),
		Member: bizId,
	}).Err()
}

func (i *InteractiveRedisCache) topLikedKey(biz string) string {
	return fmt.Sprintf("top_liked_%s", biz)
}

// BatchGet implements cache.InteractiveCache.
func (i *InteractiveRedisCache) MustBatchGet(
	ctx context.Context,
	biz string,
	bizIDs []int64,
) ([]domain.Interactive, error) {
	res := make([]domain.Interactive, 0, len(bizIDs))
	for _, bizID := range bizIDs {
		intr, err := i.Get(ctx, biz, bizID)
		if err != nil {
			return []domain.Interactive{}, err
		}
		res = append(res, intr)
	}
	return res, nil
}

// BatchSet implements cache.InteractiveCache.
func (i *InteractiveRedisCache) BatchSet(
	ctx context.Context,
	biz string,
	bizIDs []int64,
	intr []domain.Interactive,
) error {
	var err error
	for idx, bizID := range bizIDs {
		er := i.Set(ctx, biz, bizID, intr[idx])
		if er != nil {
			// log the error
			err = er
		}
	}
	return err
}

// Set implements cache.InteractiveCache.
func (i *InteractiveRedisCache) Set(
	ctx context.Context,
	biz string,
	bizID int64,
	intr domain.Interactive,
) error {
	key := i.Key(biz, bizID)
	err := i.client.HSet(
		ctx,
		key,
		fieldCollectCnt,
		intr.CollectCnt,
		fieldReadCnt,
		intr.ReadCnt,
		fieldLikeCnt,
		intr.LikeCnt,
	).Err()
	if err != nil {
		return err
	}

	return i.client.Expire(ctx, key, intrExpiryTime).Err()
}

// Get implements cache.InteractiveCache.
func (i *InteractiveRedisCache) Get(
	ctx context.Context,
	biz string,
	bizID int64,
) (domain.Interactive, error) {
	key := i.Key(biz, bizID)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		// No data
		return domain.Interactive{}, errors.New("no data")
	}
	var intr domain.Interactive
	intr.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	intr.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	intr.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	intr.Biz = biz
	intr.BizID = bizID
	return intr, nil
}

// DecrCollectCntIfPresent implements cache.InteractiveCache.
func (i *InteractiveRedisCache) DecrCollectCntIfPresent(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.Key(biz, bizID)}, fieldCollectCnt, -1).Err()
}

// IncrCollectorCntIfPresent implements cache.InteractiveCache.
func (i *InteractiveRedisCache) IncrCollectCntIfPresent(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.Key(biz, bizID)}, fieldCollectCnt, 1).Err()
}

// DecrLikeCntIfPresent implements cache.InteractiveCache.
func (i *InteractiveRedisCache) DecrLikeCntIfPresent(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.Key(biz, bizID)}, fieldLikeCnt, -1).Err()
}

// IncrLikeCntIfPresent implements cache.InteractiveCache.
func (i *InteractiveRedisCache) IncrLikeCntIfPresent(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.Key(biz, bizID)}, fieldLikeCnt, 1).Err()
}

// IncrReadCntIfPresent implements cache.InteractiveCache.
func (i *InteractiveRedisCache) IncrReadCntIfPresent(
	ctx context.Context,
	biz string,
	bizID int64,
) error {
	return i.client.Eval(ctx, luaIncrCnt, []string{i.Key(biz, bizID)}, fieldReadCnt, 1).Err()
}

func NewInteractiveRedisCache(client redis.Cmdable) cache.InteractiveCache {
	return &InteractiveRedisCache{
		client: client,
	}
}
