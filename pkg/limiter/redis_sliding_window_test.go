package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisSlidingWindowLimiter(t *testing.T) {
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
		windowsAmount   int
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
				key := "ip-limit-0.0.0.0"
				err := rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:                     context.Background(),
			prefix:                  "ip-limit",
			biz:                     "0.0.0.0",
			interval:                100 * time.Millisecond,
			windowsAmount:           10,
			limit:                   100,
			requests:                10,
			delayedRequests:         10,
			expectedLimitedRequests: 0,
		},
		// {
		// 	name: "limit",
		// 	after: func(t *testing.T) {
		// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		// 		defer cancel()
		// 		key := "ip-limit-0.0.0.0"
		// 		err := rdb.Del(ctx, key).Err()
		// 		assert.NoError(t, err)
		// 	},
		// 	ctx:                     context.Background(),
		// 	prefix:                  "ip-limit",
		// 	biz:                     "0.0.0.0",
		// 	interval:                1000 * time.Millisecond,
		// 	windowsAmount:           10,
		// 	limit:                   100,
		// 	requests:                120,
		// 	delayedRequests:         0,
		// 	expectedLimitedRequests: 2,
		// },
		// {
		// 	name: "windows amount set to 0, use default value 10",
		// 	after: func(t *testing.T) {
		// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		// 		defer cancel()
		// 		key := "ip-limit-0.0.0.0"
		// 		err := rdb.Del(ctx, key).Err()
		// 		assert.NoError(t, err)
		// 	},
		// 	ctx:                     context.Background(),
		// 	prefix:                  "ip-limit",
		// 	biz:                     "0.0.0.0",
		// 	interval:                100 * time.Millisecond,
		// 	windowsAmount:           0,
		// 	limit:                   100,
		// 	requests:                12,
		// 	delayedRequests:         10,
		// 	expectedLimitedRequests: 2,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.after(t)

			limiter := NewRedisSlidingWindowLimiter(&RedisSlidingWindowOptions{
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

			windowsAmount := tc.windowsAmount
			if tc.windowsAmount == 0 {
				windowsAmount = 10
			}

			delay := tc.interval.Nanoseconds() / int64(windowsAmount)
			time.Sleep(2 * time.Duration(delay))

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
