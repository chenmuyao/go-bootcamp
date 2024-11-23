package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicatedUser = errors.New("email already exists")
	ErrRecordNotFound = gorm.ErrRecordNotFound
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
	ID       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	Phone sql.NullString `gorm:"unique"`

	// NOTE: UTC-0
	Ctime int64
	Utime int64

	Name     string `gorm:"type=varchar(128)"`
	Birthday int64
	Profile  string `gorm:"type=varchar(4096)"`
}

func (dao *UserDAO) Insert(ctx context.Context, u User) (User, error) {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	res := dao.db.WithContext(ctx).Create(&u)
	err := res.Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr = 1062
		if me.Number == duplicateErr {
			// email conflict
			return User{}, ErrDuplicatedUser
		}
	}
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
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
