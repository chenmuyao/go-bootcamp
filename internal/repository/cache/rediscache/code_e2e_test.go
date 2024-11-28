package rediscache

import (
	"context"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCodeCacheSet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name:   "send ok",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "600123", code)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9+time.Second+50)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "12345",
			code:    "600123",
			wantErr: nil,
		},
		{
			name: "sent too many",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				err := rdb.Set(ctx, key, "600123", time.Minute*10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "600123", code)
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "12345",
			code:    "600123",
			wantErr: cache.ErrCodeSendTooMany,
		},
		{
			name: "no expiration",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				err := rdb.Set(ctx, key, "600123", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:12345"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "600123", code)
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "12345",
			code:    "600123",
			wantErr: ErrNoCodeExp,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			cc := NewCodeRedisCache(rdb)

			err := cc.Set(tc.ctx, tc.biz, tc.phone, tc.code)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
