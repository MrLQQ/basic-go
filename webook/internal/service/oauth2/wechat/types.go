package wechat

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wecaht/callback")

type service struct {
	appID     string
	appSecret string
	client    *http.Client
}

func NewService(appid string, appSecret string) Service {
	return &service{
		appID:     appid,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const baseURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	// 这是另外一种写法
	queryParams := url.Values{}
	queryParams.Set("appid", s.appID)
	queryParams.Set("secret", s.appSecret)
	queryParams.Set("code", code)
	queryParams.Set("grant_type", "authorization_code")
	accessTokenURL := baseURL + "?" + queryParams.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenURL, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var res Result
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("调用微信接口失败 errocde %d, errmsg %s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}

func (s service) AuthURL(ctx context.Context) (string, error) {
	const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprint(authURLPattern, s.appID, redirectURL, state), nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errMsg"`

	Scope string `json:"scope"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenId  string `json:"openid"`
	UnionId string `json:"unionid"`
}
