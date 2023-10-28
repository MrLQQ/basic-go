package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"fmt"
	"math/rand"
)

var ErrCodeSendTooMany = cache.ErrCodeSendTooMany

type ICodeService interface {
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
	Send(ctx context.Context, biz, phone string) error
}

// CodeRedisService ---------------------------------------使用Redis缓存实现--------------------------------------------------//
type CodeRedisService struct {
	ICodeService
	repo *repository.CodeRedisRepository
	sms  sms.Service
}

func NewCodeRedisService(repo *repository.CodeRedisRepository, smsSve sms.Service) *CodeRedisService {
	return &CodeRedisService{
		repo: repo,
		sms:  smsSve,
	}
}

func (svc *CodeRedisService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 这里是不是要开始发验证吗？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeRedisService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		// 相当于，我们对外屏蔽的验证次数过多的错误，只告诉调用者，验证错误
		return false, nil
	}
	return ok, err
}

func (svc *CodeRedisService) generate() string {
	// 0~999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

// CodeCacheService --------------------------------------------本地缓存实现-----------------------------------------------//
type CodeCacheService struct {
	ICodeService
	repo *repository.CodeCacheRepository
	sms  sms.Service
}

func NewCodeCacheService(repo *repository.CodeCacheRepository, smsSve sms.Service) *CodeCacheService {
	return &CodeCacheService{
		repo: repo,
		sms:  smsSve,
	}
}

func (svc *CodeCacheService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	// 这里是不是要开始发验证吗？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeCacheService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		// 相当于，我们对外屏蔽的验证次数过多的错误，只告诉调用者，验证错误
		return false, nil
	}
	return ok, err
}

func (svc *CodeCacheService) generate() string {
	// 0~999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
