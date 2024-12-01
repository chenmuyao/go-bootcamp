package limiter

import (
	"fmt"
	"sync"
	"time"
)

// {{{ TokenBucketLimiter

type TokenBucketOptions struct {
	Prefix string

	Capacity int

	// to calculate rate
	RelaseAmount int
	Interval     time.Duration
}

type tokenBucketRateInfo struct {
	lastReleaseTime time.Time
	tokens          int
}

type TokenBucketLimiter struct {
	cache         map[string]tokenBucketRateInfo
	prefix        string
	capacity      int
	releaseAmount int
	interval      time.Duration

	mutex sync.Mutex
}

func NewTokenBucketLimiter(options *TokenBucketOptions) *TokenBucketLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	return &TokenBucketLimiter{
		prefix:        prefix,
		cache:         map[string]tokenBucketRateInfo{},
		capacity:      options.Capacity,
		releaseAmount: options.RelaseAmount,
		interval:      options.Interval,
		mutex:         sync.Mutex{},
	}
}

func (fw *TokenBucketLimiter) release(
	now time.Time,
	rateInfo tokenBucketRateInfo,
) tokenBucketRateInfo {
	timePassed := now.Sub(rateInfo.lastReleaseTime)
	intervalsPassed := timePassed.Nanoseconds() / fw.interval.Nanoseconds()
	shouldRelease := intervalsPassed * int64(fw.releaseAmount)

	rateInfo.tokens += int(shouldRelease)

	if rateInfo.tokens >= fw.capacity {
		rateInfo.tokens = fw.capacity
	}
	rateInfo.lastReleaseTime = now
	return rateInfo
}

func (fw *TokenBucketLimiter) AcceptConnection(biz string) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	now := time.Now()

	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, ok := fw.cache[key]

	if !ok {
		// Not found
		fw.cache[key] = tokenBucketRateInfo{
			tokens:          fw.capacity - 1,
			lastReleaseTime: now,
		}
		return true
	}

	rateInfo := fw.release(now, res)
	rateInfo.tokens--

	if rateInfo.tokens < 0 {
		rateInfo.tokens = 0
		fw.cache[key] = rateInfo
		return false
	}

	fw.cache[key] = rateInfo

	return true
}

// }}}
