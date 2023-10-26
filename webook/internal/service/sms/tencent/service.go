package tencent

import (
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
	// 处理异常
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			// 基本上不可能进来这里
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Dk" {
			// 发送失败
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

func NewService(client *sms.Client, appId string, SignName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		SignName: &SignName,
	}
}
