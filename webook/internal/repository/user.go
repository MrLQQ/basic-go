package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Edit(ctx context.Context, profile domain.UserProfile) error
	Profile(ctx context.Context, profile domain.UserProfile) (domain.UserProfile, error)
	FindByWechat(ctx *gin.Context, openId string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewCacheUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomainUser(u), nil
}

func (repo *CacheUserRepository) toDomainUser(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		WechatInfo: domain.WechatInfo{
			OpenId:  u.WechatOpenId.String,
			UnionId: u.WechatUnionId.String,
		},
	}
}
func (repo *CacheUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:    u.Id,
		Email: sql.NullString{String: u.Email, Valid: u.Email != ""},
		Phone: sql.NullString{String: u.Phone, Valid: u.Phone != ""},
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionId,
			Valid:  u.WechatInfo.UnionId != "",
		},
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenId,
			Valid:  u.WechatInfo.OpenId != "",
		},
		Password: u.Password,
	}
}

func (repo *CacheUserRepository) Edit(ctx context.Context, profile domain.UserProfile) error {
	return repo.dao.Edit(ctx, dao.UserProfile{
		Id:       profile.Id,
		User_id:  profile.User_id,
		Nickname: profile.Nickname,
		Birthday: profile.Birthday,
		About_me: profile.About_me,
	})
}

func (repo *CacheUserRepository) Profile(ctx context.Context, profile domain.UserProfile) (domain.UserProfile, error) {
	du, err := repo.cache.Get(ctx, profile.User_id)
	if err == nil {
		return du, err
	}

	u, err := repo.dao.Profile(ctx, dao.UserProfile{
		User_id: profile.User_id,
	})
	if err != nil {
		return domain.UserProfile{}, err
	}
	du = repo.toDomainProfile(u)
	err = repo.cache.Set(ctx, du)
	return du, nil
}

func (repo *CacheUserRepository) toDomainProfile(u dao.UserProfile) domain.UserProfile {
	return domain.UserProfile{
		Id:       u.Id,
		User_id:  u.User_id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		About_me: u.About_me,
	}
}

func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomainUser(u), nil
}

func (repo *CacheUserRepository) FindByWechat(ctx *gin.Context, openId string) (domain.User, error) {
	ue, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomainUser(ue), nil
}
