package limiter

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// {{{ RedisFixedWindowLimiter

//go:embed lua/fixed_window.lua
var luaFixedWindow string

type RedisFixedWindowOptions struct {
	RedisClient redis.Cmdable
	Prefix      string
	Interval    time.Duration
	Limit       int
}

type RedisFixedWindowLimiter struct {
	cmd      redis.Cmdable
	prefix   string
	interval int64
	limit    int
}

func NewRedisFixedWindowLimiter(options *RedisFixedWindowOptions) *RedisFixedWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	return &RedisFixedWindowLimiter{
		cmd:      options.RedisClient,
		prefix:   prefix,
		interval: options.Interval.Milliseconds(),
		limit:    options.Limit,
	}
}

func (fw *RedisFixedWindowLimiter) AcceptConnection(ctx context.Context, biz string) bool {
	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, err := fw.cmd.Eval(ctx, luaFixedWindow, []string{key}, fw.limit, fw.interval).Int()
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
