package limiter

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

//go:embed lua/token_bucket.lua
var luaTokenBucket string

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type redisTokenBucketLimiter struct {
	cmd           redis.Cmdable
	prefix        string
	capacity      int
	releaseAmount int
	interval      int64
}

func NewRedisTokenBucketLimiter(options *RedisTokenBucketOptions) *redisTokenBucketLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	return &redisTokenBucketLimiter{
		cmd:           options.RedisClient,
		prefix:        prefix,
		interval:      options.Interval.Milliseconds(),
		capacity:      options.Capacity,
		releaseAmount: options.ReleaseAmount,
	}
}

// }}}
// {{{ Other structs

type RedisTokenBucketOptions struct {
	RedisClient redis.Cmdable
	Prefix      string

	Capacity int

	// to calculate rate
	ReleaseAmount int
	Interval      time.Duration
}

// }}}
// {{{ Struct Methods

func (fw *redisTokenBucketLimiter) AcceptConnection(ctx context.Context, biz string) bool {
	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, err := fw.cmd.Eval(ctx, luaTokenBucket, []string{key}, fw.releaseAmount, fw.interval, fw.capacity).
		Int()
	if err != nil {
		// Redis error, limit by default
		slog.Error("redis error", "err", err)
		return false
	}

	switch res {
	case -1: // limit
		return false
	default: // accept
		return true
	}
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
