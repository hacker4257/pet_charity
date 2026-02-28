package sms

import (
	"encoding/json"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AliyunSms struct {
	client       *dysmsapi.Client
	signName     string
	templateCode string
}

func NewAliyunSms(accessKeyID, accessKeySecret, signName, templateCode string) (*AliyunSms, error) {
	cfg := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	}
	client, err := dysmsapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create sms client failed:%w", err)
	}
	return &AliyunSms{
		client:       client,
		signName:     signName,
		templateCode: templateCode,
	}, nil
}

//发送验证码
func (s *AliyunSms) SendCode(phone, code string) error {
	paramMap := map[string]string{"code": code}
	paramJSON, _ := json.Marshal(paramMap)

	request := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(s.signName),
		TemplateCode:  tea.String(s.templateCode),
		TemplateParam: tea.String(string(paramJSON)),
	}

	resp, err := s.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("send sms failed: %w", err)
	}

	if *resp.Body.Code != "OK" {
		return fmt.Errorf("sms error: %s", *resp.Body.Message)
	}

	return nil
}
