package limiter

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// {{{ RedisSlidingWindowLimiter

//go:embed lua/sliding_window.lua
var luaSlidingWindow string

type RedisSlidingWindowOptions struct {
	RedisClient   redis.Cmdable
	Prefix        string
	Interval      time.Duration
	WindowsAmount int
	Limit         int
}

type RedisSlidingWindowLimiter struct {
	cmd        redis.Cmdable
	prefix     string
	windowSize int64
	limit      int
}

func NewRedisSlidingWindowLimiter(options *RedisSlidingWindowOptions) *RedisSlidingWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	windowsAmount := options.WindowsAmount
	if options.WindowsAmount == 0 {
		windowsAmount = 10
	}
	windowSize := options.Interval.Milliseconds() / int64(windowsAmount)
	return &RedisSlidingWindowLimiter{
		cmd:        options.RedisClient,
		prefix:     prefix,
		windowSize: windowSize,
		limit:      options.Limit / windowsAmount,
	}
}

func (fw *RedisSlidingWindowLimiter) AcceptConnection(ctx context.Context, biz string) bool {
	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, err := fw.cmd.Eval(ctx, luaSlidingWindow, []string{key}, fw.limit, fw.windowSize).Int()
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
