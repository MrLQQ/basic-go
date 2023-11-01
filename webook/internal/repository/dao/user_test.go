package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// lastInsertID: INSERT语句会有这个部分
				// rowAffected: 受到影响的行数，比如说插入了多少行、更新了多少行、删除了多少行
				mockRes := sqlmock.NewResult(1, 1)
				// 这里传入的是sql的正则表达式
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			user: User{
				Email: sql.NullString{String: "280@qq.com", Valid: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewGORMUserDao(db)
			err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
