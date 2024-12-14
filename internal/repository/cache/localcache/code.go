package localcache

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
)

// {{{ Consts

const (
	// if a second same request is comming in X% of expiration time, then it is
	// considered as too fraquent
	percentageTooFrequent = 10
	maxRetryTime          = 3
)

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type CodeLocalCache struct {
	cache.BaseCodeCache
	code       *ttlcache.Cache[string, string]
	cnt        *ttlcache.Cache[string, int]
	expiration time.Duration

	mu sync.Mutex
}

func NewCodeLocalCache(
	code *ttlcache.Cache[string, string],
	cnt *ttlcache.Cache[string, int],
	expiration time.Duration,
) cache.CodeCache {
	return &CodeLocalCache{
		code:       code,
		cnt:        cnt,
		expiration: expiration,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (c *CodeLocalCache) Set(ctx context.Context, biz, phone, code string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.Key(biz, phone)

	codeItem := c.code.Get(key)
	if codeItem != nil {
		slog.Info("time until", "time", time.Until(codeItem.ExpiresAt()))
	}
	// check if insert too frequently
	if codeItem != nil &&
		time.Until(codeItem.ExpiresAt()) >= c.getExpThreshold() {
		// request too frequently
		return cache.ErrCodeSendTooMany
	}
	// not too many, just replace it with a new one
	// or codeItem == nil

	c.code.Set(key, code, ttlcache.DefaultTTL)
	c.cnt.Set(key, maxRetryTime, ttlcache.DefaultTTL)

	return nil
}

func (c *CodeLocalCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.Key(biz, phone)

	codeItem := c.code.Get(key)
	codeCnt := c.cnt.Get(key)

	if codeItem == nil || codeCnt == nil {
		// no code associated
		return false, nil
	}
	if codeCnt.Value() == 0 {
		// No chances left
		return false, cache.ErrCodeVerifyTooMany
	}

	// verify
	if code != codeItem.Value() {
		newCode := codeCnt.Value()
		newCode--
		c.cnt.Set(key, newCode, codeItem.TTL())
		return false, nil
	}

	// ok
	c.cnt.Set(key, 0, codeItem.TTL())
	return true, nil
}

func (c *CodeLocalCache) getExpThreshold() time.Duration {
	return c.expiration * (100 - percentageTooFrequent) / 100
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
