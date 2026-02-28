package payment

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

type Alipay struct {
	client          *alipay.Client
	alipayPublicKey string
	notifyURL       string
	returnURL       string
}

func NewAlipay(appID, privateKeyPath, alipayPublicPath, notifyURL, returnURL string) (*Alipay, error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read alipay private key failed: %w", err)
	}

	client, err := alipay.NewClient(appID, string(privateKey), false)
	if err != nil {
		return nil, fmt.Errorf("create alipay failed: %w", err)
	}

	alipayPublickey, err := os.ReadFile(alipayPublicPath)
	if err != nil {
		return nil, fmt.Errorf("read alipay public key failed:%w", err)
	}

	client.SetNotifyUrl(notifyURL)
	client.SetReturnUrl(returnURL)

	if err := client.SetCertSnByContent(nil, nil, []byte(alipayPublickey)); err != nil {
		client.AutoVerifySign([]byte(alipayPublickey))
	}

	return &Alipay{
		client:          client,
		alipayPublicKey: string(alipayPublickey),
		notifyURL:       notifyURL,
		returnURL:       returnURL,
	}, nil
}

//创建电脑网站支付
func (a *Alipay) CreatePageOrder(orderNo, subject string, amountCent int64) (string, error) {
	//支付宝金额单位是元， 需要转换
	amountYuan := fmt.Sprintf("%.2f", float64(amountCent)/100)

	bm := make(gopay.BodyMap)
	bm.Set("subject", subject)
	bm.Set("out_trade_no", orderNo)
	bm.Set("total_amount", amountYuan)
	bm.Set("product_code", "FAST_INSTANT_TRADE_PAY")

	payURL, err := a.client.TradePagePay(context.Background(), bm)
	if err != nil {
		return "", fmt.Errorf("create alipay order failed: %w", err)
	}

	return payURL, nil

}

// VerifyNotify 解析并验证支付宝回调通知
// 验证 RSA2 签名，检查交易状态，返回商户订单号
func (a *Alipay) VerifyNotify(req *http.Request) (string, error) {
	// 1. 解析回调参数
	notifyData, err := alipay.ParseNotifyToBodyMap(req)
	if err != nil {
		return "", fmt.Errorf("parse alipay notify failed: %w", err)
	}

	// 2. 验证签名
	if _, err := alipay.VerifySign(a.alipayPublicKey, notifyData); err != nil {
		return "", fmt.Errorf("verify alipay sign failed: %w", err)
	}

	// 3. 检查交易状态
	tradeStatus := notifyData.Get("trade_status")
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		return "", nil // 非成功状态，返回空让调用方忽略
	}

	return notifyData.Get("out_trade_no"), nil
}
