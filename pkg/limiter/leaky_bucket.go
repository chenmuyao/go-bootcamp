package limiter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// {{{ LeakyBucketLimiter

type LeakyBucketOptions struct {
	Prefix string

	Capacity int

	// to calculate rate
	Limit    int
	Interval time.Duration
}

type leakyBucketRateInfo struct {
	lastLeakTime time.Time
	water        int
}

type leakyBucketLimiter struct {
	cache    map[string]leakyBucketRateInfo
	prefix   string
	capacity int
	limit    int
	interval time.Duration

	mutex sync.Mutex
}

func NewLeakyBucketLimiter(options *LeakyBucketOptions) *leakyBucketLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	return &leakyBucketLimiter{
		prefix:   prefix,
		cache:    map[string]leakyBucketRateInfo{},
		capacity: options.Capacity,
		limit:    options.Limit,
		interval: options.Interval,
		mutex:    sync.Mutex{},
	}
}

func (fw *leakyBucketLimiter) leak(
	now time.Time,
	rateInfo leakyBucketRateInfo,
) leakyBucketRateInfo {
	timePassed := now.Sub(rateInfo.lastLeakTime)
	intervalsPassed := timePassed.Nanoseconds() / fw.interval.Nanoseconds()
	shouldLeak := intervalsPassed * int64(fw.limit)

	rateInfo.water -= int(shouldLeak)

	if rateInfo.water < 0 {
		rateInfo.water = 0
	}
	rateInfo.lastLeakTime = now
	return rateInfo
}

func (fw *leakyBucketLimiter) AcceptConnection(ctx context.Context, biz string) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	now := time.Now()

	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, ok := fw.cache[key]

	if !ok {
		// Not found
		fw.cache[key] = leakyBucketRateInfo{
			water:        1,
			lastLeakTime: now,
		}
		return true
	}

	rateInfo := fw.leak(now, res)
	rateInfo.water++

	if rateInfo.water > fw.capacity {
		rateInfo.water = fw.capacity
		fw.cache[key] = rateInfo
		return false
	}

	fw.cache[key] = rateInfo

	return true
}

// }}}
