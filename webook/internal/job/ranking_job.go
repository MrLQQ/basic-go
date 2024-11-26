package job

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	l       logger.LoggerV1
	timeout time.Duration
	client  *rlock.Client
	key     string
}

func NewRankingJob(
	svc service.RankingService,
	l logger.LoggerV1,
	timeout time.Duration) *RankingJob {
	return &RankingJob{svc: svc,
		key:     "job:ranking",
		l:       l,
		timeout: timeout}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	lock, err := r.client.Lock(ctx, r.key, r.timeout,
		&rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
			// 重试的超时
		}, time.Second)
	if err != nil {
		return err
	}
	defer func() {
		// 解锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := lock.Unlock(ctx)
		if er != nil {
			r.l.Error("ranking job释放分布锁失败", logger.Error(er))
		}
	}()
	ctx, cancel = context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}
