package limiter

import (
	"fmt"
	"sync"
	"time"
)

// {{{ FixedWindowLimiter

type FixedWindowOptions struct {
	Prefix   string
	Interval time.Duration
	Limit    int
}

type fixedWindowRateInfo struct {
	timeBegin time.Time
	count     int
}

type FixedWindowLimiter struct {
	cache    map[string]fixedWindowRateInfo
	prefix   string
	interval time.Duration
	limit    int
	mutex    sync.Mutex
}

func NewFixedWindowLimiter(options *FixedWindowOptions) *FixedWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "rate-limit"
	}
	return &FixedWindowLimiter{
		prefix:   prefix,
		interval: options.Interval,
		limit:    options.Limit,
		cache:    map[string]fixedWindowRateInfo{},
		mutex:    sync.Mutex{},
	}
}

func (fw *FixedWindowLimiter) AcceptConnection(biz string) bool {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	now := time.Now()

	key := fmt.Sprintf("%s-%s", fw.prefix, biz)
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
