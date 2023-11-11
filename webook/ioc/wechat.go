package ioc

import (
	"basic-go/webook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	//appID, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("WECHAT_APP_ID is not found")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("WECHAT_APP_SECRET is not found")
	//}
	appID := "appID"
	appSecret := "appSecret"
	return wechat.NewService(appID, appSecret)
}
