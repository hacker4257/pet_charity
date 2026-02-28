package service

import (
	"errors"
	"net/http"
	"time"

	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/event"
	"github.com/hacker4257/pet_charity/pkg/logger"
	"github.com/hacker4257/pet_charity/pkg/payment"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type DonationService struct {
	donationRepo repository.DonationRepository
	orgRepo      repository.OrgRepository
	petRepo      repository.PetRepository
	wechatPay    *payment.WechatPay
	alipay       *payment.Alipay
	txManager    repository.TransactionManager
}

func NewDonationService(
	donaRepo repository.DonationRepository,
	orgRepo repository.OrgRepository,
	petRepo repository.PetRepository,
	txManager repository.TransactionManager,
) *DonationService {
	svc := &DonationService{
		donationRepo: donaRepo,
		orgRepo:      orgRepo,
		petRepo:      petRepo,
		txManager:    txManager,
	}

	//支付
	wCfg := config.Global.WechatPay
	if wCfg.MchID != "" {
		wp, err := payment.NewWechatPay(
			wCfg.MchID, wCfg.SerialNo, wCfg.APIKeyV3,
			wCfg.PrivateKeyPath, wCfg.NotifyURL,
		)
		if err != nil {
			logger.Error("init wechat pay failed", logger.Err(err))
		} else {
			svc.wechatPay = wp
		}
	}

	aCfg := config.Global.AliPay
	if aCfg.AppID != "" {
		ap, err := payment.NewAlipay(
			aCfg.AppID, aCfg.PrivateKeyPath,
			aCfg.AlipayPublicKeyPath, aCfg.NotifyURL, aCfg.ReturnURL,
		)
		if err != nil {
			logger.Errorf("init alipay failed: %v", err)
		} else {
			svc.alipay = ap
		}
	}

	return svc
}

type CreateDonationRequest struct {
	TargetType    string `json:"target_type" binding:"required,oneof=platform organization pet"`
	TargetID      uint   `json:"target_id"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Message       string `json:"message" binding:"max=500"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=wechat alipay"`
}

type CreateDonationResponse struct {
	DonationID uint   `json:"donation_id"`
	TradeNo    string `json:"trade_no"`
	PayURL     string `json:"pay_url"`
}

//创建捐赠订单
func (s *DonationService) Create(userID uint, req *CreateDonationRequest) (*CreateDonationResponse, error) {
	//1.验证捐赠对象
	if err := s.validateTarget(req.TargetType, req.TargetID); err != nil {
		return nil, err
	}

	//2.金额校验 （最小一分 = 0.01元， 最大100万 = 1亿分）
	if req.Amount > 100000000 {
		return nil, errors.New("amount too large")
	}
	//2.1 防止重复提交：同一用户对同一目标有未支付订单，直接复用
	existing, err := s.donationRepo.FindPendingByUser(userID, req.TargetType, req.TargetID)
	if err == nil {
		//已有未支付订单
		payURL, err := s.createPayURL(existing.TradeNo, existing.Amount, req.PaymentMethod)
		if err != nil {
			return nil, err
		}
		return &CreateDonationResponse{
			DonationID: existing.ID,
			TradeNo:    existing.TradeNo,
			PayURL:     payURL,
		}, nil
	}
	//3生成订单号
	tradeNo := utils.GenerateOrderNo("D")

	//4. 创建数据集记录
	donation := &model.Donation{
		UserID:        userID,
		TargetType:    req.TargetType,
		TargetID:      req.TargetID,
		Amount:        req.Amount,
		Message:       req.Message,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: "pending",
		TradeNo:       tradeNo,
	}

	if err := s.donationRepo.Create(donation); err != nil {
		return nil, errors.New("create donation failed")
	}
	s.donationRepo.AddExpireTask(tradeNo, time.Now().Add(30*time.Minute))
	//5.调用支付平台创建预支付

	payURL, err := s.createPayURL(tradeNo, req.Amount, req.PaymentMethod)
	if err != nil {
		return nil, errors.New("create payURL faield")
	}
	return &CreateDonationResponse{
		DonationID: donation.ID,
		TradeNo:    tradeNo,
		PayURL:     payURL,
	}, nil
}

func (s *DonationService) createPayURL(tradeNo string, amount int64, method string) (string, error) {
	description := "Pet Charity Donation"
	switch method {
	case "wechat":
		if s.wechatPay == nil {
			return "", errors.New("wechat pay not configured")
		}
		return s.wechatPay.CreateNativeOrder(
			config.Global.WechatPay.AppID, tradeNo, description, amount,
		)
	case "alipay":
		if s.alipay == nil {
			return "", errors.New("alipay not configured")
		}
		return s.alipay.CreatePageOrder(tradeNo, description, amount)
	}
	return "", errors.New("unsupported payment method")

}

//验证捐赠对象
func (s *DonationService) validateTarget(targetType string, targetID uint) error {
	switch targetType {
	case "platform":
		return nil
	case "organization":
		if targetID == 0 {
			return errors.New("organization id is required")
		}
		org, err := s.orgRepo.FindByID(targetID)
		if err != nil {
			return errors.New("organization not found")
		}
		if org.Status != "approved" {
			return errors.New("organization is not approved")
		}
	case "pet":
		if targetID == 0 {
			return errors.New("pet id is required")
		}
		_, err := s.petRepo.FindByID(targetID)
		if err != nil {
			return errors.New("pet not found")
		}
	}
	return nil
}

//支付回调处理
func (s *DonationService) HandlePaymentCallback(tradeNo string) error {
	//1. 查找订单
	donation, err := s.donationRepo.FindByTradeNo(tradeNo)
	if err != nil {
		return errors.New("donation not found")
	}

	//2.幂等检查：已支付的不再处理
	if donation.PaymentStatus == "paid" {
		return nil
	}
	//3. 更新订单状态
	now := time.Now()
	var affected int64
	err = s.txManager.Transaction(func(tx repository.TransactionContext) error {
		var txErr error
		affected, txErr = tx.DonationRepo().UpdateFieldsWithStatus(donation.ID, "pending", map[string]interface{}{
			"payment_status": "paid",
			"paid_at":        now,
		})
		return txErr
	})
	if err != nil {
		return err
	}
	//检查affected
	if affected == 0 {
		return nil
	}
	event.Publish("donation", donation.UserID, donation.ID)
	return nil
}

// HandleWechatNotify 处理微信支付V3回调
// 解密验证回调内容，提取订单号并更新支付状态
func (s *DonationService) HandleWechatNotify(req *http.Request) error {
	if s.wechatPay == nil {
		return errors.New("wechat pay not configured")
	}
	tradeNo, err := s.wechatPay.ParseNotify(req)
	if err != nil {
		return err
	}
	return s.HandlePaymentCallback(tradeNo)
}

// HandleAlipayNotify 处理支付宝回调
// 验证签名、检查交易状态，提取订单号并更新支付状态
func (s *DonationService) HandleAlipayNotify(req *http.Request) error {
	if s.alipay == nil {
		return errors.New("alipay not configured")
	}
	tradeNo, err := s.alipay.VerifyNotify(req)
	if err != nil {
		return err
	}
	if tradeNo == "" {
		return nil // 非成功交易状态，忽略
	}
	return s.HandlePaymentCallback(tradeNo)
}

//查询订单状态
func (s *DonationService) GetStatus(userID uint, donationID uint) (*model.Donation, error) {
	donation, err := s.donationRepo.FindByID(donationID)
	if err != nil {
		return nil, errors.New("donation not found")
	}
	if donation.UserID != userID {
		return nil, errors.New("permission denied")
	}
	return donation, nil
}

//用户捐赠记录
func (s *DonationService) ListByUser(userID uint, page, pageSize int) ([]model.Donation, int64, error) {
	return s.donationRepo.ListByUser(userID, page, pageSize)
}

//捐赠公示
func (s *DonationService) ListPublic(targetType string, targetID uint, page, pageSize int) ([]model.Donation, int64, error) {
	return s.donationRepo.ListPublic(targetType, targetID, page, pageSize)
}

func (s *DonationService) Stats() (*repository.DonationStats, error) {
	return s.donationRepo.Stats()
}

// StartExpireWorker 定时检查并关闭过期订单
func (s *DonationService) StartExpireWorker() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			tradeNos, err := s.donationRepo.GetExpiredTasks()
			if err != nil || len(tradeNos) == 0 {
				continue
			}
			for _, tradeNo := range tradeNos {
				// 用 tradeNo 关闭更准确，需要一个新方法
				s.closeByTradeNo(tradeNo)
				s.donationRepo.RemoveExpireTask(tradeNo)
			}
		}
	}()
}

func (s *DonationService) closeByTradeNo(tradeNo string) {
	donation, err := s.donationRepo.FindByTradeNo(tradeNo)
	if err != nil {
		return
	}
	if donation.PaymentStatus != "pending" {
		return
	}
	s.donationRepo.UpdateFieldsWithStatus(donation.ID, "pending", map[string]interface{}{
		"payment_status": "expired",
	})
}
