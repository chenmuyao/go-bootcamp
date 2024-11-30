package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlidingWindowLimiter(t *testing.T) {
	testCases := []struct {
		name string

		ip              string
		interval        time.Duration
		windowsAmount   int
		limit           int
		requests        int
		delayedRequests int

		expectedLimitedRequests int
	}{
		{
			name:                    "pass",
			ip:                      "0.0.0.0",
			interval:                100 * time.Millisecond,
			windowsAmount:           10,
			limit:                   100,
			requests:                10,
			delayedRequests:         10,
			expectedLimitedRequests: 0,
		},
		{
			name:                    "limited",
			ip:                      "0.0.0.0",
			interval:                100 * time.Millisecond,
			windowsAmount:           10,
			limit:                   100,
			requests:                12,
			delayedRequests:         10,
			expectedLimitedRequests: 2,
		},
		{
			name:                    "windows amount set to 0, use default value 10",
			ip:                      "0.0.0.0",
			interval:                100 * time.Millisecond,
			windowsAmount:           0,
			limit:                   100,
			requests:                12,
			delayedRequests:         10,
			expectedLimitedRequests: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			limiter := NewSlidingWindowLimiter(&SlidingWindowOptions{
				Interval:      tc.interval,
				WindowsAmount: tc.windowsAmount,
				Limit:         tc.limit,
			})
			resCh := make(chan bool)

			for range tc.requests {
				go func() {
					res := limiter.AcceptConnection(tc.ip)
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
			time.Sleep(time.Duration(delay))

			for range tc.delayedRequests {
				go func() {
					res := limiter.AcceptConnection(tc.ip)
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
