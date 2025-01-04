package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

type Interactive struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`

	// <biz_id, biz>
	BizID      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"uniqueIndex:biz_type_id,length:128"`
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Utime      int64
	Ctime      int64
}

// IncrReadCnt implements InteractiveDAO.
func (g *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	now := time.Now().UnixMilli()

	// Upsert
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("`read_cnt` + 1"), // NOTE: don't forget ``
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizID:   bizID,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}
