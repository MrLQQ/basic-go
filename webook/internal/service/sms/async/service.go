package async

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/logger"
	"context"
	"time"
)

type Service struct {
	svc  sms.Service
	repo repository.AsyncSmsRepository
	l    logger.LoggerV1
}

func (s Service) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	if s.needAsync() {
		// 需要异步发送，直接存到数据库
		err := s.repo.Add(ctx, domain.AsyncSms{
			TplId:    tplId,
			Args:     args,
			Numbers:  number,
			RetryMax: 3,
		})
		return err
	}
	return s.svc.Send(ctx, tplId, args, number...)
}

func NewService(svc sms.Service,
	repo repository.AsyncSmsRepository,
	l logger.LoggerV1) *Service {
	res := &Service{
		svc:  svc,
		repo: repo,
		l:    l,
	}
	go func() {
		res.StartAsyncCycle()
	}()
	return res
}

func (s Service) StartAsyncCycle() {
	time.Sleep(time.Second * 3)
	for {
		s.AsyncSend()
	}
}

func (s Service) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	as, err := s.repo.PreemptWaitingSMS(ctx)
	cancel()
	switch err {
	case nil:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.svc.Send(ctx, as.TplId, as.Args, as.Numbers...)
		if err != nil {
			// 啥也不干
			s.l.Error("执行异步发送短信失败",
				logger.Error(err),
				logger.Field{Key: "id", Value: as.Id})
		}
		res := err == nil
		// 通知repository 这次的执行结果
		err = s.repo.ReportScheduleResult(ctx, as.Id, res)
		if err != nil {
			s.l.Error("执行异步发送消息成功，但是数据库更新失败", logger.Error(err))
		}
	case repository.ErrWaitingSMSNotFound:
		time.Sleep(time.Second)
	default:
		// 数据库出现问题
		s.l.Error("抢占异步发送短信任务失败")
		time.Sleep(time.Second)

	}
}

func (s Service) needAsync() bool {
	return true
}
