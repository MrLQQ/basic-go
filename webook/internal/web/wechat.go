package web

import (
	"basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}
}
func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	val, err := o.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "构造跳转URL失败"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Data: val})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {

}
