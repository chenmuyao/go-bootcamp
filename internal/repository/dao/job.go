package dao

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"gorm.io/gorm"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, jid int64) error
	UpdateUtime(ctx context.Context, jid int64) error
	UpdateNextTime(ctx context.Context, jid int64, nextTime time.Time) error
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)

type Job struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"unique;type:varchar(128)"`
	Executor string
	CronExpr string
	Config   string

	Status int

	Version int

	NextTime int64 `gorm:"index"`

	Utime int64
	Ctime int64
}

type GORMJobDAO struct {
	db *gorm.DB
	l  logger.Logger
}

// UpdateNextTime implements JobDAO.
func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, jid int64, nextTime time.Time) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jid).Updates(map[string]any{
		"next_time": nextTime.UnixMilli(),
	}).Error
}

// UpdateUtime implements JobDAO.
func (g *GORMJobDAO) UpdateUtime(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jid).Updates(map[string]any{
		"utime": now,
	}).Error
}

// Preempt implements JobDAO.
func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	// NOTE: Optimist lock
	var j Job
	db := g.db.WithContext(ctx)
	for {
		now := time.Now().UnixMilli()
		// Missing finding jobs that are not refreshed to execute
		err := db.Where("status = ? AND next_time < ?", jobStatusWaiting, now).First(&j).Error
		if err != nil {
			return j, err
		}
		res := db.Model(&Job{}).
			Where("id = ? AND version = ?", j.ID, j.Version).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"version": j.Version + 1,
				"utime":   now,
			})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			continue
		}
		return j, nil
	}
}

// Release implements JobDAO.
func (g *GORMJobDAO) Release(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jid).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  now,
	}).Error
}

func NewGORMJobDAO(db *gorm.DB, l logger.Logger) JobDAO {
	return &GORMJobDAO{
		db: db,
		l:  l,
	}
}
