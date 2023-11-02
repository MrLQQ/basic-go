package failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	// 当前第几个节点
	idx int32
	// 连续几个超时了
	cnt int32
	// 切换的阈值，只读的
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{svcs: svcs, threshold: threshold}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 重置cnt次数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, number...)
	switch err {
	case nil:
		// 没有任何错误，重置计数器
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		// 超时后， 增加超时次数
		atomic.AddInt32(&t.cnt, 1)
	default:
		// 遇到了错误，但是又不是超时错误
		// 我可以增加，也可以并增加
		// 如果强调一定是超时，那么就不增加
		// 如果是EOF这类的错误，还可以考虑直接切换
	}
	return err
}
