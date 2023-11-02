package retelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/limiter"
	"context"
	"errors"
)

var ErrLimited = errors.New("触发限流")

type RateLimitSMSService struct {
	// 被装饰的
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

type RateLimitSMSServiceV1 struct {
	sms.Service
	limiter limiter.Limiter
	key     string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return ErrLimited
	}
	return r.svc.Send(ctx, tplId, args, number...)
}

func NewRateLimitSMSService(svc sms.Service, limiter limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{svc: svc, limiter: limiter, key: "sms-limiter"}
}
