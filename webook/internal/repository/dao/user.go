package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		// 冲突错误码：“Error 1062 (23000): Duplicate entry '280235109@qq.com' for key 'users.email'”
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) Edit(ctx *gin.Context, profile UserProfile) error {
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

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置唯一索引
	Email    string `gorm:"unique"`
	Password string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}

type UserProfile struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置唯一索引
	User_id  string `gorm:"unique"`
	Nickname string
	Birthday string
	About_me string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
