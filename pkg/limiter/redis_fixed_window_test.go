package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCodeCacheSet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name string

		after func(t *testing.T)

		ctx             context.Context
		prefix          string
		biz             string
		interval        time.Duration
		limit           int
		requests        int
		delayedRequests int

		expectedLimitedRequests int
	}{
		{
			name: "pass",
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "ip-limit:0.0.0.0:cnt"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
				timeBeginKey := "ip-limit:0.0.0.0:time"
				err = rdb.Del(ctx, timeBeginKey).Err()
				assert.NoError(t, err)
			},
			ctx:                     context.Background(),
			prefix:                  "ip-limit",
			biz:                     "0.0.0.0",
			interval:                100 * time.Millisecond,
			limit:                   10,
			requests:                10,
			delayedRequests:         0,
			expectedLimitedRequests: 0,
		},
		{
			name: "limit",
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "ip-limit:0.0.0.0:cnt"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
				timeBeginKey := "ip-limit:0.0.0.0:time"
				err = rdb.Del(ctx, timeBeginKey).Err()
				assert.NoError(t, err)
			},
			ctx:                     context.Background(),
			prefix:                  "ip-limit",
			biz:                     "0.0.0.0",
			interval:                100 * time.Millisecond,
			limit:                   10,
			requests:                12,
			delayedRequests:         10,
			expectedLimitedRequests: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.after(t)

			limiter := NewRedisFixedWindowLimiter(&RedisFixedWindowOptions{
				RedisClient: rdb,
				Prefix:      tc.prefix,
				Interval:    tc.interval,
				Limit:       tc.limit,
			})

			resCh := make(chan bool)

			for range tc.requests {
				go func() {
					res := limiter.AcceptConnection(tc.ctx, tc.biz)
					resCh <- res
				}()
			}

			var limited int

			for range tc.requests {
				if <-resCh == false {
					limited++
				}
			}

			time.Sleep(tc.interval)

			for range tc.delayedRequests {
				go func() {
					res := limiter.AcceptConnection(tc.ctx, tc.biz)
					resCh <- res
				}()
			}

			for range tc.delayedRequests {
				if <-resCh == false {
					limited++
				}
			}

			assert.Equal(t, tc.expectedLimitedRequests, limited)
		})
	}
}
