package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -source=./interactive.go -package=daomocks -destination=./mocks/interactive.mock.go
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizID int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIDs []int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizID int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizID int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	DeleteCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	Get(ctx context.Context, biz string, bizID int64) (Interactive, error)
	GetAll(ctx context.Context, biz string, limit int, offset int) ([]Interactive, error)
	MustBatchGet(ctx context.Context, biz string, bizIDs []int64) ([]Interactive, error)
	GetLikeInfo(ctx context.Context, biz string, bizID int64, uid int64) (UserLikeBiz, error)
	GetCollectInfo(
		ctx context.Context,
		biz string,
		bizID int64,
		uid int64,
	) (UserCollectionBiz, error)
	GetByIDs(ctx context.Context, biz string, ids []int64) ([]Interactive, error)
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

type UserCollectionBiz struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`
	// One ressource can only be put into one collection.
	// Otherwise the composite index should include CID too
	UID   int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizID int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz   string `gorm:"uniqueIndex:uid_biz_type_id,length:128"`
	// collection ID
	CID   int64 `gorm:"index"`
	Utime int64
	Ctime int64
}

// GetByIDs implements InteractiveDAO.
func (g *GORMInteractiveDAO) GetByIDs(
	ctx context.Context,
	biz string,
	ids []int64,
) ([]Interactive, error) {
	var intrs []Interactive
	err := g.db.WithContext(ctx).
		Where("biz = ? AND biz_id IN ?", biz, ids).
		Find(&intrs).
		Error
	return intrs, err
}

// GetAll implements InteractiveDAO.
func (g *GORMInteractiveDAO) GetAll(
	ctx context.Context,
	biz string,
	limit int,
	offset int,
) ([]Interactive, error) {
	var intrs []Interactive
	err := g.db.WithContext(ctx).
		Where("biz = ?", biz).
		Offset(offset).
		Limit(limit).
		Find(&intrs).
		Error
	return intrs, err
}

// BatchGet implements InteractiveDAO.
func (g *GORMInteractiveDAO) MustBatchGet(
	ctx context.Context,
	biz string,
	bizIDs []int64,
) ([]Interactive, error) {
	res := make([]Interactive, 0, len(bizIDs))
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMInteractiveDAO(tx)
		for _, bizID := range bizIDs {
			intr, err := txDAO.Get(ctx, biz, bizID)
			if err != nil {
				return err
			}
			res = append(res, intr)
		}
		return nil
	})
	if err != nil {
		return []Interactive{}, err
	}
	return res, nil
}

// Get implements InteractiveDAO.
func (g *GORMInteractiveDAO) Get(
	ctx context.Context,
	biz string,
	bizID int64,
) (Interactive, error) {
	var intr Interactive
	err := g.db.WithContext(ctx).Where("biz_id = ? AND biz = ?", bizID, biz).First(&intr).Error
	return intr, err
}

// GetLikeInfo implements InteractiveDAO.
func (g *GORMInteractiveDAO) GetLikeInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := g.db.WithContext(ctx).
		Where("uid = ? AND biz_id = ? AND biz = ? AND status = ?", uid, bizID, biz, 1).
		First(&res).
		Error
	return res, err
}

func (g *GORMInteractiveDAO) GetCollectInfo(
	ctx context.Context,
	biz string,
	bizID int64,
	uid int64,
) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := g.db.WithContext(ctx).
		Where("uid = ? AND biz_id = ? AND biz = ?", uid, bizID, biz).
		First(&res).
		Error
	return res, err
}

// DeleteCollectionBiz implements InteractiveDAO.
func (g *GORMInteractiveDAO) DeleteCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Delete(&UserCollectionBiz{}, map[string]interface{}{
			"uid":    cb.UID,
			"biz_id": cb.BizID,
			"biz":    cb.Biz,
		}).Error
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).
			Model(&Interactive{}).
			Where("biz_id = ? AND biz = ?", cb.BizID, cb.Biz).
			Updates(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` - 1"), // NOTE: don't forget ``
				"utime":       now,
			}).
			Error
	})
}

// InsertCollectionBiz implements InteractiveDAO.
func (g *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"), // NOTE: don't forget ``
				"utime":       now,
			}),
		}).Create(&Interactive{
			Biz:        cb.Biz,
			BizID:      cb.BizID,
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
		}).Error
	})
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

// BatchIncrReadCnt implements InteractiveDAO.
func (g *GORMInteractiveDAO) BatchIncrReadCnt(
	ctx context.Context,
	bizs []string,
	bizIDs []int64,
) error {
	// NOTE: Use transaction to improve the performance for batch update
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMInteractiveDAO(tx)
		for i, bizID := range bizIDs {
			err := txDAO.IncrReadCnt(ctx, bizs[i], bizID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}
