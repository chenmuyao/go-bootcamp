package limiter

import (
	"fmt"
	"sync"
	"time"
)

const (
	FixedWindow = iota
	SlidingWindow
	TokenBucket
	LeakyBucket
)

// {{{ FixedWindowLimiter

type FixedWindowOptions struct {
	Interval time.Duration
	Limit    int

	// Optional
	Prefix string
}

type fixedWindowRateInfo struct {
	count     int
	timeBegin time.Time
}

type FixedWindowLimiter struct {
	prefix   string
	interval time.Duration
	limit    int

	cache map[string]fixedWindowRateInfo

	mutex sync.Mutex
}

func NewFixedWindowLimiter(options *FixedWindowOptions) *FixedWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "IP-limit"
	}
	return &FixedWindowLimiter{
		prefix:   prefix,
		interval: options.Interval,
		limit:    options.Limit,
		cache:    map[string]fixedWindowRateInfo{},
	}
}

func (fw *FixedWindowLimiter) AcceptConnection(IP string) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	now := time.Now()

	key := fmt.Sprintf("%s-%s", fw.prefix, IP)
	res, ok := fw.cache[key]
	if !ok {
		// Not found
		fw.cache[key] = fixedWindowRateInfo{
			count:     1,
			timeBegin: now,
		}
		return true
	}

	if now.Sub(res.timeBegin) > fw.interval {
		// reset
		res.count = 1
		res.timeBegin = time.Now()
		return true
	}

	res.count++
	// compare
	if res.count > fw.limit {
		// Reached the limit
		return false
	}
	fw.cache[key] = res

	return true
}

// }}}
// {{{ SlidingWindowLimiter

type slidingWindowRateInfo struct {
	requests []time.Time
}

type SlidingWindowOptions struct {
	WindowSize time.Duration
	Limit      int

	// Optional
	Prefix string
}

type SlidingWindowLimiter struct {
	prefix     string
	windowSize time.Duration
	limit      int

	cache map[string]slidingWindowRateInfo

	mutex sync.Mutex
}

func NewSlidingWindowLimiter(options *SlidingWindowOptions) *SlidingWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "IP-limit"
	}
	return &SlidingWindowLimiter{
		prefix:     prefix,
		windowSize: options.WindowSize,
		limit:      options.Limit,
		cache:      map[string]slidingWindowRateInfo{},
	}
}

func (fw *SlidingWindowLimiter) AcceptConnection(IP string) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	now := time.Now()

	key := fmt.Sprintf("%s-%s", fw.prefix, IP)
	res, ok := fw.cache[key]
	if !ok {
		// Not found
		fw.cache[key] = slidingWindowRateInfo{
			requests: []time.Time{now},
		}
		return true
	}

	// remofe old requests
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
