package limiter

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedWindowLimiter(t *testing.T) {
	ip := "0.0.0.0"
	limiter := NewFixedWindowLimiter(&FixedWindowOptions{
		Interval: time.Millisecond,
		Limit:    10,
	})

	wg := sync.WaitGroup{}
	wg.Add(10)

	for range 10 {
		go func() {
			assert.True(t, limiter.AcceptConnection(ip))
			wg.Done()
		}()
	}
	wg.Wait()
	assert.False(t, limiter.AcceptConnection(ip))
	time.Sleep(1 * time.Millisecond)
	assert.True(t, limiter.AcceptConnection(ip))
}

func TestSlidingWindowLimiter(t *testing.T) {
	ip := "0.0.0.0"
	limiter := NewSlidingWindowLimiter(&SlidingWindowOptions{
		WindowSize: 1 * time.Millisecond,
		Limit:      10,
	})

	wg := sync.WaitGroup{}
	wg.Add(10)

	for range 10 {
		go func() {
			assert.True(t, limiter.AcceptConnection(ip))
			wg.Done()
		}()
	}
	wg.Wait()
	assert.False(t, limiter.AcceptConnection(ip))
	time.Sleep(1 * time.Millisecond)
	assert.True(t, limiter.AcceptConnection(ip))
}
