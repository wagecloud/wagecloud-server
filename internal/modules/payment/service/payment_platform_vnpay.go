package paymentsvc

import (
	"context"
	"errors"
	"strconv"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/client/vnpay"
)

const (
	// TODO: move to config
	PaymentResolvePath = "/payment-resolve"
)

type VnpayPlatform struct {
	client vnpay.Client
}

func NewVnpayPlatform(client vnpay.Client) *VnpayPlatform {
	return &VnpayPlatform{
		client: client,
	}
}

func (p *VnpayPlatform) CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error) {
	return p.client.CreateOrder(ctx, vnpay.CreateOrderParams{
		PaymentID: params.PaymentID,
		Amount:    params.Amount.Float64(),
		Info:      params.Info,
		ReturnUrl: config.GetConfig().App.FrontendUrl + PaymentResolvePath,
	})
}

func (p *VnpayPlatform) VerifyPayment(ctx context.Context, data map[string]any) (string, error) {
	if err := p.client.VerifyPayment(ctx, data); err != nil {
		return "", err
	}

	tmnCode, ok := data["vnp_TmnCode"].(string)
	if !ok {
		return "", errors.New("TmnCode not found in query parameters")
	}

	// Verify the TmnCode
	if tmnCode != config.GetConfig().Vnpay.TmnCode {
		return "", errors.New("invalid TmnCode")
	}

	txnRef, ok := data["vnp_TxnRef"].(string)
	if !ok {
		return "", errors.New("transaction reference not found in query parameters")
	}

	paymentID, err := strconv.ParseInt(txnRef, 10, 64)
	if err != nil {
		return "", errors.New("invalid transaction reference")
	}

	return strconv.FormatInt(paymentID, 10), nil
}
