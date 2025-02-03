package limiter

import "context"

const (
	FixedWindow = iota
	SlidingWindow
	TokenBucket
	LeakyBucket
)

//go:generate mockgen -source=./types.go -package=limitermocks -destination=./mocks/limiter.mock.go
type Limiter interface {
	AcceptConnection(ctx context.Context, biz string) bool
}

func NewLimiter(options any) Limiter {
	switch opt := options.(type) {
	case *FixedWindowOptions:
		return NewFixedWindowLimiter(opt)
	case *SlidingWindowOptions:
		return NewSlidingWindowLimiter(opt)
	case *LeakyBucketOptions:
		return NewLeakyBucketLimiter(opt)
	case *TokenBucketOptions:
		return NewTokenBucketLimiter(opt)
	case *RedisFixedWindowOptions:
		return NewRedisFixedWindowLimiter(opt)
	case *RedisSlidingWindowOptions:
		return NewRedisSlidingWindowLimiter(opt)
	case *RedisLeakyBucketOptions:
		return NewRedisLeakyBucketLimiter(opt)
	case *RedisTokenBucketOptions:
		return NewRedisTokenBucketLimiter(opt)
	default:
		return nil
	}
}
