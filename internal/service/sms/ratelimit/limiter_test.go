package ratelimit

import (
	"context"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	smsmocks "github.com/chenmuyao/go-bootcamp/internal/service/sms/mocks"
	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	limitermocks "github.com/chenmuyao/go-bootcamp/pkg/limiter/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRateLimitSMSServiceSend(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)
		wantErr error
	}{
		{
			name: "pass",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().AcceptConnection(gomock.Any(), gomock.Any()).Return(true)
				svc.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return svc, l
			},
		},
		{
			name: "limit",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)
				l.EXPECT().AcceptConnection(gomock.Any(), gomock.Any()).Return(false)
				return svc, l
			},
			wantErr: ErrLimited,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvc, l := tc.mock(ctrl)

			svc := NewRateLimitSMSService(smsSvc, l)
			err := svc.Send(context.Background(), "12345", "text")
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
