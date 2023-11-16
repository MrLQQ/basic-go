package localsms

import (
	"basic-go/webook/pkg/logger"
	"context"
	"log"
)

type Service struct {
	l logger.LoggerV1
}

func NewService(l logger.LoggerV1) *Service {
	return &Service{l: l}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("验证码是", args)
	s.l.Info("验证码是", logger.Field{Key: "args", Value: args})
	return nil
}
