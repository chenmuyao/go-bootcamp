package cronjoblearn

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			t.Log("finished")
			return
		case now := <-ticker.C:
			t.Log("1 second", now.UnixMilli())
		}
	}
}
