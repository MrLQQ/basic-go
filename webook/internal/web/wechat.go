package web

import (
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/oauth2/wechat"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	key             []byte
	stateCookieName string
	l               logger.LoggerV1
}

func NewOAuth2WechatHandler(svc wechat.Service, hdl ijwt.Handler, userSvc service.UserService, l logger.LoggerV1) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("jBxoQWRS5L9vYr$mYq5U9d5BRPHfSBA1"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
		l:               l,
	}
}
func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	state := uuid.New()
	val, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		o.l.Error("构造跳转URL失败", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "构造跳转URL失败"})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		o.l.Error("服务异常", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "服务器异常"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Data: val})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	err := o.verifyState(ctx)
	if err != nil {
		o.l.Error("非法请求", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "非法请求"})
		return
	}
	code := ctx.Query("code")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		o.l.Warn("授权码有误", logger.Field{Key: "currentCode", Value: code}, logger.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "授权码有误"})
		return
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		o.l.Error("系统错误", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		o.l.Error("系统错误", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Code: 1, Msg: "ok"})
	return
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	cookie, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		o.l.Error("无法获得cookie", logger.Error(err))
		return fmt.Errorf("无法获得cookie %w", err)
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(cookie, &sc, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		o.l.Error("解析token失败", logger.Error(err))
		return fmt.Errorf("解析token失败 %w", err)
	}
	if state != sc.State {
		// state不匹配
		o.l.Error("state不匹配",
			logger.Field{Key: "currentState", Value: state},
			logger.Field{Key: state, Value: sc.State},
			logger.Error(err))
		return fmt.Errorf("state不匹配")
	}
	return nil
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {
		o.l.Error("系统错误", logger.Error(err))
		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
