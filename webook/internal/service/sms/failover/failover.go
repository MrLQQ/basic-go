package failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type FailOverSMSService struct {
	svcs []sms.Service

	// 这是V1的字段
	// 记录一个当前下标
	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

func (f FailOverSMSService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, number...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("发送失败，所有服务商都尝试过了")
}

// SendV1 起始下标轮询
// 并且出错也轮询
func (f FailOverSMSService) SendV1(ctx context.Context, tplId string, args []string, number ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	// 需要迭代length
	for i := idx; i < idx+length; i++ {
		// 取余来计算下标
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, number...)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return err
			// 前者是被取消，后者是超时
		}
		// 其他情况会走到这里，打印日志
		log.Println(err)
	}
	return errors.New("发送失败，所有服务商都尝试过了")
}
