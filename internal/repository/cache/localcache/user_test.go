package localcache

import (
	"context"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
	"github.com/stretchr/testify/assert"
)

func TestUserLocalCache(t *testing.T) {
	timeout := 10 * time.Millisecond
	cc := ttlcache.New(ttlcache.WithTTL[string, domain.User](timeout))
	go cc.Start()

	uc := NewUserLocalCache(cc)

	testUser := domain.User{
		ID: 123,
	}

	err := uc.Set(context.Background(), testUser)
	assert.NoError(t, err)
	got, err := uc.Get(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, testUser, got)

	_, err = uc.Get(context.Background(), 321)
	assert.ErrorIs(t, err, cache.ErrKeyNotExist)

	time.Sleep(10 * time.Millisecond)

	_, err = uc.Get(context.Background(), 123)
	assert.ErrorIs(t, err, cache.ErrKeyNotExist)
}
