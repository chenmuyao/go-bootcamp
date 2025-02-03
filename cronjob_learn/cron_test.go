package cronjoblearn

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestCronExpr(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	id, err := expr.AddFunc("@every 1s", func() {
		t.Log("1s")
	})
	assert.NoError(t, err)
	t.Log("task", id)
	expr.Start()
	time.Sleep(10 * time.Second)
	ctx := expr.Stop() // stop running new tasks
	t.Log("send stop signal")
	<-ctx.Done() // wait for the last task to finish
	t.Log("terminated")
}
