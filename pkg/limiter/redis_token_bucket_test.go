package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisTokenBucketLimiter(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name string

		after func(t *testing.T)

		ctx           context.Context
		prefix        string
		biz           string
		interval      time.Duration
		releaseAmount int
		capacity      int

		sleep           time.Duration
		requests        int
		delayedRequests int

		expectedLimitedRequests int
	}{
		{
			name: "pass",
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "ip-limit-0.0.0.0"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
			},
			ctx:           context.Background(),
			prefix:        "ip-limit",
			biz:           "0.0.0.0",
			interval:      100 * time.Millisecond,
			capacity:      100,
			releaseAmount: 10,

			sleep:                   100 * time.Millisecond,
			requests:                100,
			delayedRequests:         10,
			expectedLimitedRequests: 0,
		},
		{
			name: "limit",
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "ip-limit-0.0.0.0"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
			},
			ctx:           context.Background(),
			prefix:        "ip-limit",
			biz:           "0.0.0.0",
			interval:      100 * time.Millisecond,
			capacity:      100,
			releaseAmount: 10,

			sleep:                   100 * time.Millisecond,
			requests:                100,
			delayedRequests:         12,
			expectedLimitedRequests: 2,
		},
		{
			name: "run out of water",
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "ip-limit-0.0.0.0"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
			},
			ctx:           context.Background(),
			prefix:        "ip-limit",
			biz:           "0.0.0.0",
			interval:      10 * time.Millisecond,
			capacity:      100,
			releaseAmount: 10,

			sleep:                   120 * time.Millisecond,
			requests:                100,
			delayedRequests:         12,
			expectedLimitedRequests: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.after(t)

			limiter := NewRedisTokenBucketLimiter(&RedisTokenBucketOptions{
				RedisClient:   rdb,
				Prefix:        tc.prefix,
				Capacity:      tc.capacity,
				Interval:      tc.interval,
				ReleaseAmount: tc.releaseAmount,
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

			time.Sleep(tc.sleep)

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
