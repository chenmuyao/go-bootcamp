package repository

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"gorm.io/gorm"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

//go:generate mockgen -source=./sms.go -package=repomocks -destination=./mocks/sms.mock.go
type AsyncSMSRepository interface {
	AddSMS(ctx context.Context, toNb string, body string, args string) error
	TrySend(
		ctx context.Context,
		f func(ctx context.Context, toNb string, body string, args ...string) error,
		maxRetry int,
	)
}

// }}}
// {{{ Struct

type asyncSMSRepository struct {
	dao dao.AsyncSMSDAO
	db  *gorm.DB
}

func NewAsyncSMSRepository(dao dao.AsyncSMSDAO, db *gorm.DB) AsyncSMSRepository {
	return &asyncSMSRepository{
		dao: dao,
		db:  db,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (r *asyncSMSRepository) AddSMS(
	ctx context.Context,
	toNb string,
	body string,
	args string,
) error {
	return r.dao.Insert(ctx, dao.SMSInfo{
		ToNb: toNb,
		Body: body,
		Args: args,
	})
}

func (r *asyncSMSRepository) TrySend(
	ctx context.Context,
	f func(ctx context.Context, toNb string, body string, args ...string) error,
	maxRetry int,
) {
	for {
		err := r.db.Transaction(func(tx *gorm.DB) error {
			sms := dao.NewAsyncSMSDAO(tx)

			return trySendFirstRecord(ctx, sms, f, maxRetry)
		})
		if err != nil {
			return
		}
	}
}

// }}}
// {{{ Private functions

func trySendFirstRecord(
	ctx context.Context,
	smsDAO dao.AsyncSMSDAO,
	f func(ctx context.Context, toNb string, body string, args ...string) error,
	maxRetry int,
) error {
	// Get the first record
	first, err := smsDAO.GetFirst(ctx)
	if err != nil {
		// If nothing is found, we will return error here, and stop the loop
		return err
	}

	var args []string
	if err := json.Unmarshal([]byte(first.Args), &args); err != nil {
		// should not happen
		slog.Error("json unmarshal error", "err", err)
		_ = smsDAO.Delete(ctx, first)
		return nil
	}

	// Try to send it
	err = f(ctx, first.ToNb, first.Body, args...)
	if err == dao.ErrDuplicatedSMS {
		// If failed for the same reasons, stop and return
		if first.RetryTimes+1 >= maxRetry {
			// Delete the request
			// NOTE: don't handle error here, we don't care if it is really deleted actually
			_ = smsDAO.Delete(ctx, first)
			return nil
		}
		first.RetryTimes++
		// Update retry times for this request and try the next one
		// NOTE: don't handle update error here. The sending error is more important,
		// we can't do anything anyway.
		_ = smsDAO.Update(ctx, first)
		return nil
	}
	// If success, delete the record
	// and try sending the next one, until running into an error or stop
	_ = smsDAO.Delete(ctx, first)
	return nil
}

// }}}
// {{{ Package functions

// }}}
