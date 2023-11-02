package failover

import (
	"basic-go/webook/internal/service/sms"
	smsmocks "basic-go/webook/internal/service/sms/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTimeoutFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) []sms.Service
		threshold int32
		idx       int32
		cnt       int32

		wantErr error
		wantIdx int32
		wantCnt int32
	}{
		{
			name: "没有触发切换",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
			idx:       0,
			cnt:       12,
			threshold: 15,
			wantIdx:   0,
			// 成功了 重置超时次数
			wantCnt: 0,
			wantErr: nil,
		},
		{
			name: "触发切换，发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			// 触发了切换
			wantIdx: 1,
			wantCnt: 0,
			wantErr: nil,
		},
		{
			name: "触发切换，失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			wantIdx:   1,
			wantCnt:   0,
			wantErr:   errors.New("发送失败"),
		},
		{
			name: "触发切换，超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			wantIdx:   1,
			wantCnt:   1,
			wantErr:   context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewTimeoutFailoverSMSService(tc.mock(ctrl), tc.threshold)
			svc.idx = tc.idx
			svc.cnt = tc.cnt
			err := svc.Send(context.Background(), "123", []string{}, "11234123")
			assert.Equal(t, tc.wantErr, err)
			t.Log("tc.wantIdx=", tc.wantIdx)
			t.Log("svc.idx=", svc.idx)
			assert.Equal(t, tc.wantIdx, svc.idx)
			t.Log("tc.wantCnt=", tc.wantCnt)
			t.Log("svc.cnt=", svc.cnt)
			assert.Equal(t, tc.wantCnt, svc.cnt)
		})
	}
}
