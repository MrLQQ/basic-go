package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type jwtHandler struct {
}

func (h *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 30分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTkey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.JSON(http.StatusOK, Result{Code: 0, Msg: "登录成功"})
}

var JWTkey = []byte("jBxoQWRS5L9vYr$mYq5U9d5BRPHfSBAe")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
