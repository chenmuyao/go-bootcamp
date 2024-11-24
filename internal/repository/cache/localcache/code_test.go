package localcache

import (
	"context"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
	"github.com/stretchr/testify/assert"
)

func TestLocalCache(t *testing.T) {
	timeout := 100 * time.Millisecond
	code := ttlcache.New(
		ttlcache.WithTTL[string, string](timeout),
		ttlcache.WithDisableTouchOnHit[string, string](),
	)
	cnt := ttlcache.New(
		ttlcache.WithTTL[string, int](timeout),
		ttlcache.WithDisableTouchOnHit[string, int](),
	)
	go code.Start()
	go cnt.Start()

	uc := NewCodeLocalCache(code, cnt, timeout)
	ctx := context.Background()

	err := uc.Set(ctx, "login", "12345", "123456")
	assert.NoError(t, err)
	err = uc.Set(ctx, "login", "12345", "123456")
	assert.ErrorIs(t, cache.ErrCodeSendTooMany, err)
	time.Sleep(10 * time.Millisecond)
	err = uc.Set(ctx, "login", "12345", "123456")
	assert.NoError(t, err)

	ok, err := uc.Verify(ctx, "login", "123456", "123456")
	assert.False(t, ok)
	assert.NoError(t, err)

	for range 3 {
		ok, err = uc.Verify(ctx, "login", "12345", "654321")
		assert.False(t, ok)
		assert.NoError(t, err)
	}
	ok, err = uc.Verify(ctx, "login", "12345", "654321")
	assert.False(t, ok)
	assert.ErrorIs(t, cache.ErrCodeVerifyTooMany, err)

	time.Sleep(10 * time.Millisecond)
	err = uc.Set(ctx, "login", "12345", "123456")
	assert.NoError(t, err)

	ok, err = uc.Verify(ctx, "login", "12345", "123456")
	assert.True(t, ok)
	assert.NoError(t, err)
}
