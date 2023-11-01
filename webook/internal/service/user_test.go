package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	repomocks "basic-go/webook/internal/repository/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123456#password")
	encryptd, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encryptd))
	err = bcrypt.CompareHashAndPassword(encryptd, []byte("wrong password"))
	assert.NotNil(t, err)
}

func Test_userService_Login(t *testing.T) {
	testCase := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1234@qq.com").AnyTimes().
					Return(domain.User{
						Email: "1234@qq.com",
						// 加密后的正确密码
						Password: "$2a$10$.cEhcQ/KJzv0b7BJfh.MW..n16W78HREJD7jS852YB3Eb46jZyGyC",
						Phone:    "1764040400044",
					}, nil)
				return repo
			},
			email: "1234@qq.com",
			// 用户输入没有加密的
			password: "123456#password",

			wantUser: domain.User{
				Email:    "1234@qq.com",
				Password: "$2a$10$.cEhcQ/KJzv0b7BJfh.MW..n16W78HREJD7jS852YB3Eb46jZyGyC",
				Phone:    "1764040400044",
			},
			wantErr: nil,
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1234@qq.com").AnyTimes().
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email: "1234@qq.com",
			// 用户输入没有加密的
			password: "123456#password",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewuserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			t.Log(tc.wantErr)
			t.Log(err)
			t.Log(tc.wantUser)
			t.Log(user)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
