package domain

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Job struct {
	ID         int64
	Name       string
	Config     string
	CronExpr   string
	Executor   string
	CancelFunc func()
}

func (j Job) NextTime() time.Time {
	c := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	s, _ := c.Parse(j.CronExpr)
	return s.Next(time.Now())
}
