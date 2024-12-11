package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	daomocks "github.com/chenmuyao/go-bootcamp/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestAsyncSMSRepository_TrySend(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) dao.AsyncSMSDAO

		// inputs
		sendFunc func(ctx context.Context, toNb string, body string, args ...string) error
		maxRetry int

		wantErr error
	}{
		{
			name: "pass",
			mock: func(ctrl *gomock.Controller) dao.AsyncSMSDAO {
				smsDAO := daomocks.NewMockAsyncSMSDAO(ctrl)

				smsDAO.EXPECT().GetFirst(gomock.Any()).Return(dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"]`,
					RetryTimes: 0,
				}, nil)

				smsDAO.EXPECT().Delete(gomock.Any(), dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"]`,
					RetryTimes: 0,
				}).Return(nil)
				return smsDAO
			},
			sendFunc: func(ctx context.Context, toNb, body string, args ...string) error {
				return nil
			},
		},
		{
			name: "get from db error",
			mock: func(ctrl *gomock.Controller) dao.AsyncSMSDAO {
				smsDAO := daomocks.NewMockAsyncSMSDAO(ctrl)

				smsDAO.EXPECT().GetFirst(gomock.Any()).Return(dao.SMSInfo{}, errors.New("db error"))
				return smsDAO
			},
			sendFunc: func(ctx context.Context, toNb, body string, args ...string) error {
				return nil
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "json unmarshall error",
			mock: func(ctrl *gomock.Controller) dao.AsyncSMSDAO {
				smsDAO := daomocks.NewMockAsyncSMSDAO(ctrl)

				smsDAO.EXPECT().GetFirst(gomock.Any()).Return(dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"`,
					RetryTimes: 0,
				}, nil)

				smsDAO.EXPECT().Delete(gomock.Any(), dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"`,
					RetryTimes: 0,
				}).Return(nil)
				return smsDAO
			},
			sendFunc: func(ctx context.Context, toNb, body string, args ...string) error {
				return nil
			},
		},
		{
			name: "retry failed for the first time",
			mock: func(ctrl *gomock.Controller) dao.AsyncSMSDAO {
				smsDAO := daomocks.NewMockAsyncSMSDAO(ctrl)

				smsDAO.EXPECT().GetFirst(gomock.Any()).Return(dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"]`,
					RetryTimes: 0,
				}, nil)

				smsDAO.EXPECT().Update(gomock.Any(), dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"]`,
					RetryTimes: 1,
				}).Return(nil)
				return smsDAO
			},
			sendFunc: func(ctx context.Context, toNb, body string, args ...string) error {
				return dao.ErrDuplicatedSMS
			},
			maxRetry: 3,
		},
		{
			name: "reached max retry",
			mock: func(ctrl *gomock.Controller) dao.AsyncSMSDAO {
				smsDAO := daomocks.NewMockAsyncSMSDAO(ctrl)

				smsDAO.EXPECT().GetFirst(gomock.Any()).Return(dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"]`,
					RetryTimes: 2,
				}, nil)

				smsDAO.EXPECT().Delete(gomock.Any(), dao.SMSInfo{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.UnixMilli(0),
						UpdatedAt: time.UnixMilli(0),
					},
					ToNb:       "12345",
					Body:       "body",
					Args:       `["arg1","arg2"]`,
					RetryTimes: 2,
				}).Return(nil)
				return smsDAO
			},
			sendFunc: func(ctx context.Context, toNb, body string, args ...string) error {
				return dao.ErrDuplicatedSMS
			},
			maxRetry: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			smsDAO := tc.mock(ctrl)

			err := trySendFirstRecord(context.Background(), smsDAO, tc.sendFunc, tc.maxRetry)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
