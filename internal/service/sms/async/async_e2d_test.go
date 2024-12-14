package async

import (
	"context"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/config"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	smsmocks "github.com/chenmuyao/go-bootcamp/internal/service/sms/mocks"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/ratelimit"
	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAsyncSMSCode(t *testing.T) {
	config.InitConfig("../../../../config/dev.yaml")
	db, err := gorm.Open(
		mysqlDriver.Open(config.Cfg.DB.DSN),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		t.Fatal("failed to init tables")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		mock func(ctrl *gomock.Controller) sms.Service

		// rate limit
		limit    int
		interval time.Duration

		// async send
		retryTimes           int
		retryErrorCodes      []int
		pollInterval         time.Duration
		goroutineRunningTime time.Duration

		wantErr error
	}{
		{
			name:   "no rate limit",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "rate-limit-sms-limiter:cnt"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
				timeKey := "rate-limit-sms-limiter:time"
				err = rdb.Del(ctx, timeKey).Err()
				assert.NoError(t, err)
			},
			mock: func(ctrl *gomock.Controller) sms.Service {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsSvc.EXPECT().Send(gomock.Any(), "to", "body", "arg1", "arg2").Return(nil)
				return smsSvc
			},
			limit:    10,
			interval: 1 * time.Second,
		},
		{
			name:   "no rate limit, hit server error, resend ok",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "rate-limit-sms-limiter:cnt"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
				timeKey := "rate-limit-sms-limiter:time"
				err = rdb.Del(ctx, timeKey).Err()
				assert.NoError(t, err)

				// should be no record in db
				var s dao.SMSInfo
				err = db.WithContext(ctx).First(&s).Error
				assert.ErrorContains(t, err, "record not found")
			},
			mock: func(ctrl *gomock.Controller) sms.Service {
				smsSvc := smsmocks.NewMockService(ctrl)
				gomock.InOrder(
					smsSvc.EXPECT().
						Send(gomock.Any(), "to", "body", "arg1", "arg2").
						Return(&sms.SMSThirdPartyError{
							Code: 123,
							Msg:  "Internal server error",
						}),
					smsSvc.EXPECT().Send(gomock.Any(), "to", "body", "arg1", "arg2").Return(nil),
				)
				return smsSvc
			},
			limit:                10,
			interval:             1 * time.Second,
			retryTimes:           3,
			retryErrorCodes:      []int{123},
			pollInterval:         10 * time.Millisecond,
			goroutineRunningTime: 100 * time.Millisecond,
		},
		{
			name: "rate limit, resend ok",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "rate-limit-sms-limiter:cnt"
				err := rdb.Set(ctx, cntKey, 1, 1*time.Second).Err()
				assert.NoError(t, err)

				timeKey := "rate-limit-sms-limiter:time"
				err = rdb.Set(ctx, timeKey, time.Now().UnixMilli(), 1*time.Second).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "rate-limit-sms-limiter:cnt"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
				timeKey := "rate-limit-sms-limiter:time"
				err = rdb.Del(ctx, timeKey).Err()
				assert.NoError(t, err)

				// should be no record in db
				var s dao.SMSInfo
				err = db.WithContext(ctx).First(&s).Error
				assert.ErrorContains(t, err, "record not found")
			},
			mock: func(ctrl *gomock.Controller) sms.Service {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsSvc.EXPECT().Send(gomock.Any(), "to", "body", "arg1", "arg2").Return(nil)
				return smsSvc
			},
			limit:                1,
			interval:             10 * time.Millisecond,
			retryTimes:           3,
			retryErrorCodes:      []int{123},
			pollInterval:         20 * time.Millisecond,
			goroutineRunningTime: 100 * time.Millisecond,
		},
		{
			name: "hit error, reach max retry, abandon",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				cntKey := "rate-limit-sms-limiter:cnt"
				err := rdb.Del(ctx, cntKey).Err()
				assert.NoError(t, err)
				timeKey := "rate-limit-sms-limiter:time"
				err = rdb.Del(ctx, timeKey).Err()
				assert.NoError(t, err)

				// should be no record in db
				var s dao.SMSInfo
				err = db.WithContext(ctx).First(&s).Error
				assert.ErrorContains(t, err, "record not found")
			},
			mock: func(ctrl *gomock.Controller) sms.Service {
				smsSvc := smsmocks.NewMockService(ctrl)
				gomock.InOrder(
					// first fail
					smsSvc.EXPECT().
						Send(gomock.Any(), "to", "body", "arg1", "arg2").
						Return(&sms.SMSThirdPartyError{
							Code: 123,
							Msg:  "Internal server error",
						}),
					// 3 retries
					smsSvc.EXPECT().
						Send(gomock.Any(), "to", "body", "arg1", "arg2").
						Return(&sms.SMSThirdPartyError{
							Code: 123,
							Msg:  "Internal server error",
						}),
					smsSvc.EXPECT().
						Send(gomock.Any(), "to", "body", "arg1", "arg2").
						Return(&sms.SMSThirdPartyError{
							Code: 123,
							Msg:  "Internal server error",
						}),
					smsSvc.EXPECT().
						Send(gomock.Any(), "to", "body", "arg1", "arg2").
						Return(&sms.SMSThirdPartyError{
							Code: 123,
							Msg:  "Internal server error",
						}),
				)
				return smsSvc
			},
			limit:                10,
			interval:             1 * time.Second,
			retryTimes:           3,
			retryErrorCodes:      []int{123},
			pollInterval:         10 * time.Millisecond,
			goroutineRunningTime: 100 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			ctx, cancel := context.WithCancel(context.Background())
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc := tc.mock(ctrl)

			daoSMS := dao.NewAsyncSMSDAO(db)
			repoSMS := repository.NewAsyncSMSRepository(daoSMS, db)

			rateLimitSMSSvc := ratelimit.NewRateLimitSMSService(
				smsSvc,
				limiter.NewLimiter(&limiter.RedisFixedWindowOptions{
					RedisClient: rdb,
					Interval:    tc.interval,
					Limit:       tc.limit,
				}),
			)

			asyncSvc := NewAsyncSMSService(ctx, rateLimitSMSSvc, repoSMS, &AsyncSMSServiceOptions{
				PollInterval:    tc.pollInterval,
				RetryTimes:      tc.retryTimes,
				RetryErrorCodes: tc.retryErrorCodes,
			})

			err := asyncSvc.Send(context.Background(), "to", "body", "arg1", "arg2")
			assert.Equal(t, tc.wantErr, err)
			time.Sleep(tc.goroutineRunningTime)
			cancel()
			time.Sleep(100 * time.Millisecond)
		})
	}
}
