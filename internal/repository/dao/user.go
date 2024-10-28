package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	// NOTE: autoIncrement for performance:
	// 1. rows physically stored in key order
	// 2. Read-ahead
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	// NOTE: UTC-0
	Ctime int64
	Utime int64
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}
