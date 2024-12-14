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

//go:embed lua/leaky_bucket.lua
var luaLeakyBucket string

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type redisLeakyBucketLimiter struct {
	cmd      redis.Cmdable
	prefix   string
	capacity int
	limit    int
	interval int64
}

func NewRedisLeakyBucketLimiter(options *RedisLeakyBucketOptions) *redisLeakyBucketLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	return &redisLeakyBucketLimiter{
		cmd:      options.RedisClient,
		prefix:   prefix,
		interval: options.Interval.Milliseconds(),
		capacity: options.Capacity,
		limit:    options.Limit,
	}
}

// }}}
// {{{ Other structs

type RedisLeakyBucketOptions struct {
	RedisClient redis.Cmdable
	Prefix      string

	Capacity int

	// to calculate rate
	Limit    int
	Interval time.Duration
}

// }}}
// {{{ Struct Methods

func (fw *redisLeakyBucketLimiter) AcceptConnection(ctx context.Context, biz string) bool {
	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, err := fw.cmd.Eval(ctx, luaLeakyBucket, []string{key}, fw.limit, fw.interval, fw.capacity).
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
