package payment

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

type WechatPay struct {
	client    *wechat.ClientV3
	apiKeyV3  string
	notifyURL string
}

func NewWechatPay(mchID, serialNo, apiKeyv3, privateKeyPath, notifyURL string) (*WechatPay, error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read wechat private key failed: %w", err)
	}
	client, err := wechat.NewClientV3(mchID, serialNo, apiKeyv3, string(privateKey))
	if err != nil {
		return nil, fmt.Errorf("create wechat client failed: %w", err)
	}

	return &WechatPay{
		client:    client,
		apiKeyV3:  apiKeyv3,
		notifyURL: notifyURL,
	}, nil
}

//创建native 支付 pc端扫码支付
func (w *WechatPay) CreateNativeOrder(appID, orderNo, description string, amountCent int64) (string, error) {
	bm := make(gopay.BodyMap)
	bm.Set("appid", appID)
	bm.Set("description", description)
	bm.Set("out_trade_no", orderNo)
	bm.Set("notify_url", w.notifyURL)
	bm.SetBodyMap("amount", func(b gopay.BodyMap) {
		b.Set("total", amountCent)
		b.Set("currency", "CNY")
	})
	resp, err := w.client.V3TransactionNative(context.Background(), bm)
	if err != nil {
		return "", fmt.Errorf("create wechat order failed: %w", err)
	}
	return resp.Response.CodeUrl, nil
}

// ParseNotify 解析并验证微信支付V3回调通知
// 通过 AEAD-AES-256-GCM 解密通知内容，同时验证数据完整性
func (w *WechatPay) ParseNotify(req *http.Request) (string, error) {
	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		return "", fmt.Errorf("parse wechat notify failed: %w", err)
	}

	// 解密通知内容，AEAD 解密同时验证数据完整性 DecryptCipherText
	result, err := notifyReq.DecryptPayCipherText(w.apiKeyV3)
	if err != nil {
		return "", fmt.Errorf("decrypt wechat notify failed: %w", err)
	}

	if result.TradeState != "SUCCESS" {
		return "", fmt.Errorf("trade not success: %s", result.TradeState)
	}

	return result.OutTradeNo, nil
}
