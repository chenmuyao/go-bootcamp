package async

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	repomocks "github.com/chenmuyao/go-bootcamp/internal/repository/mocks"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	smsmocks "github.com/chenmuyao/go-bootcamp/internal/service/sms/mocks"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/ratelimit"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAsyncSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller, ctx context.Context, retryTimes int) (sms.Service, repository.AsyncSMSRepository)

		// inputs
		retryTimes           int
		retryErrorCodes      []int
		pollInterval         time.Duration
		goroutineRunningTime time.Duration

		// outputs
		wantErr error
	}{
		{
			name: "pass",
			mock: func(ctrl *gomock.Controller, ctx context.Context, retryTimes int) (sms.Service, repository.AsyncSMSRepository) {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsRepo := repomocks.NewMockAsyncSMSRepository(ctrl)

				smsSvc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				smsRepo.EXPECT().
					TrySend(ctx, gomock.Any(), retryTimes).MinTimes(1)

				return smsSvc, smsRepo
			},
			pollInterval:         10 * time.Millisecond,
			retryTimes:           3,
			goroutineRunningTime: 10 * time.Millisecond,
		},
		{
			name: "failed for system errors",
			mock: func(ctrl *gomock.Controller, ctx context.Context, retryTimes int) (sms.Service, repository.AsyncSMSRepository) {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsRepo := repomocks.NewMockAsyncSMSRepository(ctrl)

				smsSvc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("system error"))
				smsRepo.EXPECT().
					TrySend(ctx, gomock.Any(), retryTimes).MinTimes(1)
				return smsSvc, smsRepo
			},
			pollInterval:         10 * time.Millisecond,
			retryTimes:           3,
			goroutineRunningTime: 10 * time.Millisecond,
			wantErr:              errors.New("system error"),
		},
		{
			name: "failed for rate limit, add to db success",
			mock: func(ctrl *gomock.Controller, ctx context.Context, retryTimes int) (sms.Service, repository.AsyncSMSRepository) {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsRepo := repomocks.NewMockAsyncSMSRepository(ctrl)

				smsSvc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ratelimit.ErrLimited)
				smsRepo.EXPECT().
					AddSMS(gomock.Any(), "to", "body", `["arg1","arg2"]`).
					Return(nil)
				smsRepo.EXPECT().
					TrySend(ctx, gomock.Any(), retryTimes).MinTimes(1)
				return smsSvc, smsRepo
			},
			pollInterval:         10 * time.Millisecond,
			retryTimes:           3,
			goroutineRunningTime: 10 * time.Millisecond,
		},
		{
			name: "failed for rate limit, add to db failed",
			mock: func(ctrl *gomock.Controller, ctx context.Context, retryTimes int) (sms.Service, repository.AsyncSMSRepository) {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsRepo := repomocks.NewMockAsyncSMSRepository(ctrl)

				smsSvc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ratelimit.ErrLimited)
				smsRepo.EXPECT().
					AddSMS(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
				smsRepo.EXPECT().
					TrySend(ctx, gomock.Any(), retryTimes).MinTimes(1)
				return smsSvc, smsRepo
			},
			pollInterval:         10 * time.Millisecond,
			retryTimes:           3,
			goroutineRunningTime: 10 * time.Millisecond,
			wantErr:              errors.New("db error"),
		},
		{
			name: "failed for third-party client-side errors",
			mock: func(ctrl *gomock.Controller, ctx context.Context, retryTimes int) (sms.Service, repository.AsyncSMSRepository) {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsRepo := repomocks.NewMockAsyncSMSRepository(ctrl)

				smsSvc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&sms.SMSThirdPartyError{
						Code: 123,
						Msg:  "Wrong API info",
					})
				smsRepo.EXPECT().
					TrySend(ctx, gomock.Any(), retryTimes).MinTimes(1)
				return smsSvc, smsRepo
			},
			pollInterval:         10 * time.Millisecond,
			retryTimes:           3,
			goroutineRunningTime: 10 * time.Millisecond,
			wantErr: &sms.SMSThirdPartyError{
				Code: 123,
				Msg:  "Wrong API info",
			},
		},
		{
			name: "failed for third-party server-side errors, add to db success",
			mock: func(ctrl *gomock.Controller, ctx context.Context, retryTimes int) (sms.Service, repository.AsyncSMSRepository) {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsRepo := repomocks.NewMockAsyncSMSRepository(ctrl)

				smsSvc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&sms.SMSThirdPartyError{
						Code: 123,
						Msg:  "Internal server error",
					})
				smsRepo.EXPECT().
					AddSMS(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				smsRepo.EXPECT().
					TrySend(ctx, gomock.Any(), retryTimes).MinTimes(1)
				return smsSvc, smsRepo
			},
			pollInterval:         10 * time.Millisecond,
			retryTimes:           3,
			retryErrorCodes:      []int{123},
			goroutineRunningTime: 10 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, repo := tc.mock(ctrl, ctx, tc.retryTimes)

			asyncSvc := NewAsyncSMSService(ctx, smsSvc, repo, &AsyncSMSServiceOptions{
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
