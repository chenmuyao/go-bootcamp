package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicatedEmail = errors.New("email already exists")
	ErrRecordNotFound  = gorm.ErrRecordNotFound
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

	Name     string `gorm:"size:50"`
	Birthday int64
	Profile  string `gorm:"size:2000"`
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr = 1062
		if me.Number == duplicateErr {
			// email conflict
			return ErrDuplicatedEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateProfile(ctx context.Context, user User) error {
	err := dao.db.WithContext(ctx).Where("id=?", user.ID).Updates(user).Error

	return err
}
