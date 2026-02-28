package service

import (
	"errors"
	"fmt"

	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/logger"
	"github.com/hacker4257/pet_charity/pkg/sms"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type SmsService struct {
	smsRepo repository.SmsRepository
	sender  *sms.AliyunSms
}

func NewSmsService(smsRepo repository.SmsRepository) *SmsService {
	cfg := config.Global.Aliyun

	sender, err := sms.NewAliyunSms(
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
		cfg.SmsSignName,
		cfg.SmsTemplateCode,
	)
	if err != nil {
		logger.Warnf("SMS sender init failed: %v, SMS features will be unavailable", err)
	}

	return &SmsService{
		smsRepo: smsRepo,
		sender:  sender,
	}

}

type SendCodeRequest struct {
	Phone   string `json:"phone" binding:"required,len=11"`
	Purpose string `json:"purpose" binding:"required,oneof=login register"`
}

type SmsLoginRequest struct {
	Phone string `json:"phone" binding:"required,len=11"`
	Code  string `json:"code" binding:"required,len=6"`
}

//发送验证码
func (s *SmsService) SendCode(req *SendCodeRequest) error {
	//1. 检查60秒限流
	throttled, err := s.smsRepo.CheckThrottle(req.Phone)
	if err != nil {
		return errors.New("system error")
	}

	if throttled {
		return errors.New("please wait 60 seconds before retry")
	}

	//2. 检查每日上限
	count, err := s.smsRepo.IncrDaily(req.Phone)
	if err != nil {
		return errors.New("system error")
	}
	if count > 10 {
		return errors.New("daily sms limit exceeded")
	}

	//3.生成验证码
	code := utils.RandomCode(6)

	//4.存入redis
	if err := s.smsRepo.SaveCode(req.Purpose, req.Phone, code); err != nil {
		return errors.New("save code failed")
	}

	//5.设置限流
	if err := s.smsRepo.SetThrottle(req.Phone); err != nil {
		return errors.New("system error")
	}

	//6.发送短信
	if s.sender == nil {
		return errors.New("SMS service is not available")
	}
	if err := s.sender.SendCode(req.Phone, code); err != nil {
		return fmt.Errorf("send sms failed : %w", err)
	}

	return nil
}

//验证验证码
func (s *SmsService) VerifyCode(phone, code, purpose string) error {
	savedCode, err := s.smsRepo.GetCode(purpose, phone)
	if err != nil {
		return errors.New("code expired or not found")
	}

	if savedCode != code {
		return errors.New("incorrect code")
	}
	//验证成功，删除验证码
	s.smsRepo.DeleteCode(purpose, phone)
	return nil
}
