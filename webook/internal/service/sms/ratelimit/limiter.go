package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/limiter"
	"basic-go/webook/pkg/logger"
	"context"
	"errors"
)

var ErrLimited = errors.New("触发限流")

type RateLimitSMSService struct {
	// 被装饰的
	svc     sms.Service
	limiter limiter.Limiter
	key     string
	l       logger.LoggerV1
}

type RateLimitSMSServiceV1 struct {
	sms.Service
	limiter limiter.Limiter
	key     string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		r.l.Error("限流器发生错误", logger.Error(err))
		return err
	}
	if limited {
		r.l.Error("触发限流")
		return ErrLimited
	}
	return r.svc.Send(ctx, tplId, args, number...)
}

func NewRateLimitSMSService(svc sms.Service, limiter limiter.Limiter, l logger.LoggerV1) *RateLimitSMSService {
	return &RateLimitSMSService{svc: svc, limiter: limiter, key: "sms-limiter", l: l}
}
