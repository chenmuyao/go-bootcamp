package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLeakyBucketLimiter(t *testing.T) {
	testCases := []struct {
		name string

		ip       string
		interval time.Duration
		limit    int
		capacity int

		sleep           time.Duration
		requests        int
		delayedRequests int

		expectedLimitedRequests int
	}{
		{
			name:                    "pass",
			ip:                      "0.0.0.0",
			interval:                100 * time.Millisecond,
			capacity:                100,
			limit:                   10,
			sleep:                   100 * time.Millisecond,
			requests:                100,
			delayedRequests:         10,
			expectedLimitedRequests: 0,
		},
		{
			name:                    "limited",
			ip:                      "0.0.0.0",
			interval:                100 * time.Millisecond,
			capacity:                100,
			limit:                   10,
			sleep:                   100 * time.Millisecond,
			requests:                100,
			delayedRequests:         12,
			expectedLimitedRequests: 2,
		},
		{
			name:                    "run out of water",
			ip:                      "0.0.0.0",
			interval:                10 * time.Millisecond,
			capacity:                100,
			sleep:                   120 * time.Millisecond,
			limit:                   10,
			requests:                100,
			delayedRequests:         12,
			expectedLimitedRequests: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			limiter := NewLeakyBucketLimiter(&LeakyBucketOptions{
				Capacity: tc.capacity,
				Interval: tc.interval,
				Limit:    tc.limit,
			})
			resCh := make(chan bool)

			for range tc.requests {
				go func() {
					res := limiter.AcceptConnection(context.Background(), tc.ip)
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
					res := limiter.AcceptConnection(context.Background(), tc.ip)
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
