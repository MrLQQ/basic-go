package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type jwtHandler struct {
	signingMethod jwt.SigningMethod
	refreshKey    []byte
}

func newJwtHandler() jwtHandler {
	return jwtHandler{
		signingMethod: jwt.SigningMethodHS512,
		refreshKey:    []byte("jBxoQWRS5L9vYr$mYq5U9d5BRPHfSBA9"),
	}
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
	token := jwt.NewWithClaims(h.signingMethod, uc)
	tokenStr, err := token.SignedString(JWTkey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.JSON(http.StatusOK, Result{Code: 0, Msg: "登录成功"})
}

func (h *jwtHandler) setRefreshJWTToken(ctx *gin.Context, uid int64) error {
	rc := RefreshClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 七天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, rc)
	tokenStr, err := token.SignedString(h.refreshKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

// ExtractToken 根据约定 token在Authorization头部
// Authorization中的结构为：Bearer XXXX
func ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return authCode
	}
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

var JWTkey = []byte("jBxoQWRS5L9vYr$mYq5U9d5BRPHfSBAe")

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid int64
}
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
