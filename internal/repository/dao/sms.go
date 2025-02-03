package dao

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

var ErrDuplicatedSMS = errors.New("sms already exists")

// }}}
// {{{ Interface

//go:generate mockgen -source=./sms.go -package=daomocks -destination=./mocks/sms.mock.go
type AsyncSMSDAO interface {
	Insert(ctx context.Context, s SMSInfo) error
	Update(ctx context.Context, s SMSInfo) error
	GetFirst(ctx context.Context) (SMSInfo, error)
	Delete(ctx context.Context, smsInfo SMSInfo) error
}

// }}}
// {{{ Struct

type GORMAsyncSMSDAO struct {
	db *gorm.DB
}

func NewAsyncSMSDAO(db *gorm.DB) AsyncSMSDAO {
	return &GORMAsyncSMSDAO{
		db: db,
	}
}

// }}}
// {{{ Other structs

type SMSInfo struct {
	gorm.Model

	ToNb string `gorm:"index:sms,unique"`
	Body string `gorm:"index:sms,unique"`
	Args string // JSON array

	RetryTimes int
}

// }}}
// {{{ Struct Methods

func (dao *GORMAsyncSMSDAO) Insert(ctx context.Context, s SMSInfo) error {
	err := dao.db.WithContext(ctx).Create(&s).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr = 1062
		if me.Number == duplicateErr {
			// email conflict
			return ErrDuplicatedSMS
		}
	}
	return err
}

func (dao *GORMAsyncSMSDAO) Update(ctx context.Context, s SMSInfo) error {
	err := dao.db.WithContext(ctx).Save(&s).Error
	return err
}

func (dao *GORMAsyncSMSDAO) GetFirst(ctx context.Context) (SMSInfo, error) {
	var s SMSInfo
	err := dao.db.WithContext(ctx).First(&s).Error
	return s, err
}

func (dao *GORMAsyncSMSDAO) Delete(ctx context.Context, smsInfo SMSInfo) error {
	// NOTE: Hard delete because we could readd the same message later
	err := dao.db.Unscoped().Delete(&smsInfo).Error
	return err
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
