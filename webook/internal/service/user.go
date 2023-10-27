package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUser         = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户名或密码不存在")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Edit(ctx context.Context, userProfile domain.UserProfile) error {
	return svc.repo.Edit(ctx, userProfile)
}

func (svc *UserService) Profile(ctx context.Context, userProfile domain.UserProfile) (domain.UserProfile, error) {
	return svc.repo.Profile(ctx, userProfile)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先查找，是否存在
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		// 有两种情况：
		//		1. err == nil,u是可用的
		// 		2. err != nil,系统错误
		return u, err
	}
	// 用户没找到
	err = svc.repo.Create(ctx, domain.User{Phone: phone})
	// 两种可能:
	// 		1.一种是err恰好是唯一索引冲突（phone）
	//		2.err != nil, 系统错误
	if err != nil && !errors.Is(err, repository.ErrDuplicateUser) {
		return domain.User{}, err
	}

	// 要么 err == nil,要么ErrDuplicateUser，也代表用户存在
	// 由于主从延迟,刚插入数据库的内容可能查询不到,
	return svc.repo.FindByPhone(ctx, phone)
}
