package wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wecaht/callback")

type service struct {
	appID string
}

func NewService(appid string) Service {
	return &service{
		appID: appid,
	}
}

func (s service) AuthURL(ctx context.Context) (string, error) {
	const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprint(authURLPattern, s.appID, redirectURL, state), nil
}
