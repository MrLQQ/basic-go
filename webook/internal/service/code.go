package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"fmt"
	"math/rand"
)

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 这里是不是要开始发验证吗？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeService) verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		// 相当于，我们对外屏蔽的验证次数过多的错误，只告诉调用者，验证错误
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generate() string {
	// 0~999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
