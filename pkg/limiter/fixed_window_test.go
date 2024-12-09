package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedWindowLimiter(t *testing.T) {
	testCases := []struct {
		name string

		ip              string
		interval        time.Duration
		limit           int
		requests        int
		delayedRequests int

		expectedLimitedRequests int
	}{
		{
			name:                    "pass",
			ip:                      "0.0.0.0",
			interval:                100 * time.Millisecond,
			limit:                   10,
			requests:                10,
			delayedRequests:         0,
			expectedLimitedRequests: 0,
		},
		{
			name:                    "limited",
			ip:                      "0.0.0.0",
			interval:                100 * time.Millisecond,
			limit:                   10,
			requests:                12,
			delayedRequests:         10,
			expectedLimitedRequests: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			limiter := NewFixedWindowLimiter(&FixedWindowOptions{
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

			time.Sleep(tc.interval)

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
