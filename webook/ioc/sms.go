package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/localsms"
	"basic-go/webook/internal/service/sms/tencent"
	"basic-go/webook/pkg/logger"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSMSService(l logger.LoggerV1) sms.Service {
	//拥有限流的短信发送服务
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter(InitRedis(), time.Second, 100))
	// 腾讯云消息服务
	//return initTencentSMSService()
	return localsms.NewService(l)
}

func initTencentSMSService(l logger.LoggerV1) sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯SMS的secret id")
	}

	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯SMS的secret key")
	}

	client, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap_nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(client, "appid everything", "数字签名", l)
}
