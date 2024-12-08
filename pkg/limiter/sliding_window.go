package limiter

import (
	"fmt"
	"sync"
	"time"
)

// {{{ SlidingWindowLimiter

type slidingWindowRateInfo struct {
	requests []time.Time
}

type SlidingWindowOptions struct {
	Prefix        string
	Interval      time.Duration
	WindowsAmount int
	Limit         int
}

type SlidingWindowLimiter struct {
	cache      map[string]slidingWindowRateInfo
	prefix     string
	windowSize time.Duration
	limit      int
	mutex      sync.Mutex
}

func NewSlidingWindowLimiter(options *SlidingWindowOptions) *SlidingWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	windowsAmount := options.WindowsAmount
	if options.WindowsAmount == 0 {
		windowsAmount = 10
	}
	windowSize := options.Interval.Nanoseconds() / int64(windowsAmount)
	return &SlidingWindowLimiter{
		prefix:     prefix,
		windowSize: time.Duration(windowSize),
		limit:      options.Limit / windowsAmount,
		cache:      map[string]slidingWindowRateInfo{},
		mutex:      sync.Mutex{},
	}
}

func (fw *SlidingWindowLimiter) AcceptConnection(biz string) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	now := time.Now()

	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
	res, ok := fw.cache[key]
	if !ok {
		// Not found
		fw.cache[key] = slidingWindowRateInfo{
			requests: []time.Time{now},
		}
		return true
	}

	// remove old requests
	cutTime := now.Add(-fw.windowSize)

	for len(res.requests) > 0 && res.requests[0].Before(cutTime) {
		res.requests = res.requests[1:]
	}

	// check len
	if len(res.requests) >= fw.limit {
		// reached the limit
		return false
	}

	res.requests = append(res.requests, now)

	fw.cache[key] = res

	return true
}

// }}}
