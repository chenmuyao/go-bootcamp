package limiter

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// {{{ RedisLeakyBucketLimiter

//go:embed lua/leaky_bucket.lua
var luaLeakyBucket string

type RedisLeakyBucketOptions struct {
	RedisClient redis.Cmdable
	Prefix      string

	Capacity int

	// to calculate rate
	limit    int
	Interval time.Duration
}

type RedisLeakyBucketLimiter struct {
	cmd      redis.Cmdable
	prefix   string
	capacity int
	limit    int
	interval int64
}

func NewRedisLeakyBucketLimiter(options *RedisLeakyBucketOptions) *RedisLeakyBucketLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	return &RedisLeakyBucketLimiter{
		cmd:      options.RedisClient,
		prefix:   prefix,
		interval: options.Interval.Milliseconds(),
		capacity: options.Capacity,
		limit:    options.limit,
	}
}

func (fw *RedisLeakyBucketLimiter) AcceptConnection(ctx context.Context, biz string) bool {
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
