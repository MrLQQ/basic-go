package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"database/sql"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomainUser(u), nil
}

func (repo *UserRepository) toDomainUser(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
	}
}
func (repo *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Phone:    sql.NullString{String: u.Phone, Valid: u.Phone != ""},
		Password: u.Password,
	}
}

func (repo *UserRepository) Edit(ctx context.Context, profile domain.UserProfile) error {
	return repo.dao.Edit(ctx, dao.UserProfile{
		Id:       profile.Id,
		User_id:  profile.User_id,
		Nickname: profile.Nickname,
		Birthday: profile.Birthday,
		About_me: profile.About_me,
	})
}

func (repo *UserRepository) Profile(ctx context.Context, profile domain.UserProfile) (domain.UserProfile, error) {
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

func (repo *UserRepository) toDomainProfile(u dao.UserProfile) domain.UserProfile {
	return domain.UserProfile{
		Id:       u.Id,
		User_id:  u.User_id,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		About_me: u.About_me,
	}
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomainUser(u), nil
}
