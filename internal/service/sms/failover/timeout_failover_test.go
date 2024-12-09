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

func TestTimeoutFailOverSMSServiceSend(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		idx       int32
		cnt       int32
		threshold int32

		wantIdx int32
		wantCnt int32
		wantErr error
	}{
		{
			name: "pass directly",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0}
			},
			idx:       0,
			cnt:       9,
			threshold: 10,
			wantIdx:   0,
			wantCnt:   0,
			wantErr:   nil,
		},
		{
			name: "switch and time out",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       10,
			threshold: 10,
			wantIdx:   1,
			wantCnt:   1,
			wantErr:   context.DeadlineExceeded,
		},
		{
			name: "switch and passed",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       10,
			threshold: 10,
			wantIdx:   1,
			wantCnt:   0,
			wantErr:   nil,
		},
		{
			name: "switch and received other error",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().
					Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("other error"))
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       10,
			threshold: 10,
			wantIdx:   1,
			wantCnt:   0,
			wantErr:   errors.New("other error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsSvcs := tc.mock(ctrl)

			svc := NewTimeoutFailOverSMSService(smsSvcs, tc.threshold)
			svc.cnt = tc.cnt
			svc.idx = tc.idx

			err := svc.Send(context.Background(), "12345", "text")

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantIdx, svc.idx)
			assert.Equal(t, tc.wantCnt, svc.cnt)
		})
	}
}
