package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizID int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizID int64, uid int64) error
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

type UserLikeBiz struct {
	ID     int64  `gorm:"primaryKey,autoIncrement"`
	UID    int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizID  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"uniqueIndex:uid_biz_type_id,length:128"`
	Status int
	Utime  int64
	Ctime  int64
}

// DeleteLikeInfo implements InteractiveDAO.
func (g *GORMInteractiveDAO) DeleteLikeInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()

		err := tx.WithContext(ctx).
			Model(&UserLikeBiz{}).
			Where("uid = ? AND biz_id = ? AND biz = ?", uid, bizID, biz).
			Updates(map[string]interface{}{"status": 0, "utime": now}).
			Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).
			Model(&Interactive{}).
			Where("biz_id = ? AND biz = ?", bizID, biz).
			Updates(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` - 1"), // NOTE: don't forget ``
				"utime":    now,
			}).
			Error
	})
}

// InsertLikeInfo implements InteractiveDAO.
func (g *GORMInteractiveDAO) InsertLikeInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()

		err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeBiz{
			UID:    uid,
			BizID:  bizID,
			Biz:    biz,
			Status: 1,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` + 1"), // NOTE: don't forget ``
				"utime":    now,
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizID:   bizID,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
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
