package rediscache

import (
	"context"
	"strconv"
	"time"

	_ "embed"

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

type InteractiveRedisCache struct {
	cache.BaseInteractiveCache
	client redis.Cmdable
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
		return domain.Interactive{}, nil
	}
	if len(res) == 0 {
		// No data
		return domain.Interactive{}, nil
	}
	var intr domain.Interactive
	intr.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	intr.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	intr.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
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
