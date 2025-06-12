package paymentsvc

import (
	"context"

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

func (p *VnpayPlatform) CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error) {
	return p.client.CreateOrder(ctx, vnpay.CreateOrderParams{
		PaymentID: params.PaymentID,
		Amount:    params.Amount.Float64(),
		Info:      params.Info,
		ReturnUrl: config.GetConfig().App.FrontendUrl + PaymentResolvePath,
	})
}
