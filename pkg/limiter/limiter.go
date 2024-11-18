package limiter

import (
	"fmt"
	"time"
)

// const (
// 	FixedWindow = iota
// 	SlidingWindow
// 	TokenBucket
// 	LeakyBucket
// )

type Options struct {
	Interval time.Duration
	Limit    int

	// Optional
	Prefix string
}

type FixedWindowLimiter struct {
	Prefix   string
	Interval time.Duration
	Limit    int

	Cache map[string]rateInfo
}

type rateInfo struct {
	Count     int
	TimeBegin time.Time
}

func NewFixedWindowLimiter(options *Options) *FixedWindowLimiter {
	prefix := options.Prefix
	if len(prefix) == 0 {
		prefix = "IP-limit"
	}
	return &FixedWindowLimiter{
		Prefix:   prefix,
		Interval: options.Interval,
		Limit:    options.Limit,
		Cache:    map[string]rateInfo{},
	}
}

func (fw *FixedWindowLimiter) KeyPrefix(prefix string) *FixedWindowLimiter {
	fw.Prefix = prefix
	return fw
}

func (fw *FixedWindowLimiter) AcceptConnection(IP string) bool {
	key := fw.generateKey(IP)
	res, ok := fw.Cache[key]
	if !ok {
		// Not found
		fw.Cache[key] = rateInfo{
			Count:     1,
			TimeBegin: time.Now(),
		}
		return true
	}

	now := time.Now()
	if now.Sub(res.TimeBegin) <= fw.Interval {
		res.Count++
		// compare
		if res.Count <= fw.Limit {
			fw.Cache[key] = res
			return true
		} else {
			// Reached the limit
			return false
		}
	} else {
		// reset
		res.Count = 1
		res.TimeBegin = time.Now()
	}
	return true
}

func (fw *FixedWindowLimiter) generateKey(IP string) string {
	return fmt.Sprintf("%s-%s", fw.Prefix, IP)
}
