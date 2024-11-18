package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedWindowLimiter(t *testing.T) {
	ip := "0.0.0.0"
	limiter := NewFixedWindowLimiter(&Options{Interval: time.Millisecond, Limit: 10})
	for range 10 {
		assert.True(t, limiter.AcceptConnection(ip))
	}
	assert.False(t, limiter.AcceptConnection(ip))
	time.Sleep(1 * time.Millisecond)
	assert.True(t, limiter.AcceptConnection(ip))
}
