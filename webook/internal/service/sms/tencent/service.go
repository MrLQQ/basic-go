package tencent

import (
	"basic-go/webook/pkg/logger"
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

type Service struct {
	client   *sms.Client
	appId    *string
	SignName *string
	l        logger.LoggerV1
}

func (s Service) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.SignName
	request.TemplateId = ekit.ToPtr[string](tplId)
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(number)
	response, err := s.client.SendSms(request)
	s.l.Debug("请求腾讯SendSMS接口", logger.Field{Key: "req", Value: request}, logger.Field{Key: "resp", Value: response})
	// 处理异常
	if err != nil {
		s.l.Error("An API error has returned", logger.Error(err))
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			// 基本上不可能进来这里
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			// 发送失败
			s.l.Error("短息发送失败", logger.Field{Key: "code", Value: *status.Code})
			return fmt.Errorf("短息发送失败 code:%s, msg:%s", *status.Code, *status.Message)

		}
	}
	return nil

}

func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}

func NewService(client *sms.Client, appId string, SignName string, l logger.LoggerV1) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		SignName: &SignName,
		l:        l,
	}
}
