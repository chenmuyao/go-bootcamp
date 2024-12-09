package failover

import (
	"context"
	"errors"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	smsmocks "github.com/chenmuyao/go-bootcamp/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFailOverSMSServiceSend(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		wantErr error
	}{
		{
			name: "pass at the first time",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0}
			},
		},
		{
			name: "pass at the Nth time",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("some error"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0, svc1}
			},
		},
		{
			name: "failed all",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("some error"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("some error"))
				return []sms.Service{svc0, svc1}
			},
			wantErr: errFailOverAll,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvcs := tc.mock(ctrl)

			svc := NewFailOverSMSService(smsSvcs)
			err := svc.Send(context.Background(), "12345", "text")
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
