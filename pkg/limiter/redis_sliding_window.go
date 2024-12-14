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

//go:embed lua/sliding_window.lua
var luaSlidingWindow string

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type redisSlidingWindowLimiter struct {
	cmd        redis.Cmdable
	prefix     string
	windowSize int64
	limit      int
}

func NewRedisSlidingWindowLimiter(options *RedisSlidingWindowOptions) *redisSlidingWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	windowsAmount := options.WindowsAmount
	if options.WindowsAmount == 0 {
		windowsAmount = 10
	}
	windowSize := options.Interval.Nanoseconds() / int64(windowsAmount)
	return &redisSlidingWindowLimiter{
		cmd:        options.RedisClient,
		prefix:     prefix,
		windowSize: windowSize,
		limit:      options.Limit / windowsAmount,
	}
}

// }}}
// {{{ Other structs

type RedisSlidingWindowOptions struct {
	RedisClient   redis.Cmdable
	Prefix        string
	Interval      time.Duration
	WindowsAmount int
	Limit         int
}

// }}}
// {{{ Struct Methods

func (fw *redisSlidingWindowLimiter) AcceptConnection(ctx context.Context, biz string) bool {
	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, err := fw.cmd.Eval(ctx, luaSlidingWindow, []string{key}, fw.limit, fw.windowSize, time.Now().UnixNano()).
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
