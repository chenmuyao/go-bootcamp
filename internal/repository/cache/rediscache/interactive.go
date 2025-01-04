package rediscache

import (
	"context"

	_ "embed"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

const fieldReadCnt = "read_cnt"

//go:embed lua/incr_cnt.lua
var luaIncrCnt string

type InteractiveRedisCache struct {
	cache.BaseInteractiveCache
	client redis.Cmdable
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
