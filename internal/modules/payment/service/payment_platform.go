package paymentsvc

import (
	"context"

	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
)

type Platform string

const (
	PlatformVNPAY Platform = "vnpay"
	PlatformMOMO  Platform = "momo"
)

type CreateOrderParams struct {
	PaymentID int64
	Info      string
	Amount    commonmodel.Concurrency
}

// PaymentPlatform is an interface for payment platform
type PaymentPlatform interface {
	CreateOrder(ctx context.Context, params CreateOrderParams) (url string, err error)
}
