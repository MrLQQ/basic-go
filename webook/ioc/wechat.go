package ioc

import (
	"basic-go/webook/internal/service/oauth2/wechat"
	"os"
)

func InitWechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID is not found")
	}

	return wechat.NewService(appID)
}
