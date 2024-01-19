package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateUser  = errors.New("用户冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	Edit(ctx context.Context, profile UserProfile) error
	Profile(ctx context.Context, userprofile UserProfile) (UserProfile, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
}

type GORMUserDao struct {
	db *gorm.DB
}

func (dao *GORMUserDao) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wecaht_open_id=?", openId).First(&u).Error
	return u, err
}

func NewGORMUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{
		db: db,
	}
}

func (dao *GORMUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		// 冲突错误码：“Error 1062 (23000): Duplicate entry '280235109@qq.com' for key 'users.email'”
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateUser
		}
	}
	return err
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) Edit(ctx context.Context, profile UserProfile) error {
	now := time.Now().UnixMilli()
	profile.Ctime = now
	profile.Utime = now
	user := profile
	// 查询目标userid的profile是否存在
	err := dao.db.WithContext(ctx).Where("User_ID=?", profile.User_id).First(&user).Error
	// 不存在插入，存在更新
	if err != nil {
		// 插入
		err := dao.db.WithContext(ctx).Create(&profile).Error
		return err
	} else {
		// 更新
		err := dao.db.Debug().WithContext(ctx).Where("User_ID=?", profile.User_id).Updates(&profile).Error
		return err
	}
}

func (dao *GORMUserDao) Profile(ctx context.Context, userprofile UserProfile) (UserProfile, error) {
	err := dao.db.WithContext(ctx).Where("user_id=?", userprofile.User_id).First(&userprofile).Error
	return userprofile, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "phone = ?", phone).Error
	return u, err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// gorm:"unique"  设置唯一索引
	// sql.NullString 表示可以为null的列
	Email sql.NullString `gorm:"unique"`
	Phone sql.NullString `gorm:"unique"`

	// 1、如果查询要求同时使用openid和unionid 就要使用联合唯一索引
	// 2、如果查询只用openid，那么就在openid上创建索引，或者<openid,unionid>联合索引
	// 3、如果查询只有uninid，那么就在uninid上创建索引，或者<unionid,openid>联合索引
	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString
	Password      string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}

type UserProfile struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置唯一索引
	User_id  int64 `gorm:"unique"`
	Nickname string
	Birthday string
	About_me string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
