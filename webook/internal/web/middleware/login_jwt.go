package middleware

import (
	"basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		// 根据约定 token在Authorization头部
		// Authorization中的结构为：Bearer XXXX
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			// 没登陆，没有token，Authorization这个头部都没有
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			// 没登陆，Authorization中的内容是假的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTkey, nil
		})
		if err != nil {
			// token不对，token时伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// token解析出来了，但是token可能是非法的，或者过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expiresTime := uc.ExpiresAt
		if expiresTime.Before(time.Now()) {
			// 过期了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// （为了测试当前设置token过期时间为1分钟）若剩余过期时间<50s就刷新
		if expiresTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTkey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				// 这边不用中断，因为仅仅是过期时间没有刷新，但是用户是登录的
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
