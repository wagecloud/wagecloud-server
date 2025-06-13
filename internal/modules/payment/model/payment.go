package paymentmodel

import (
	"time"

	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
)

type PaymentStatus string

const (
	PaymentStatusUnknown  PaymentStatus = "PAYMENT_STATUS_UNKNOWN"
	PaymentStatusPending  PaymentStatus = "PAYMENT_STATUS_PENDING"
	PaymentStatusSuccess  PaymentStatus = "PAYMENT_STATUS_SUCCESS"
	PaymentStatusCanceled PaymentStatus = "PAYMENT_STATUS_CANCELED"
	PaymentStatusFailed   PaymentStatus = "PAYMENT_STATUS_FAILED"
)

type PaymentMethod string

const (
	PaymentMethodUnknown PaymentMethod = "PAYMENT_METHOD_UNKNOWN"
	PaymentMethodVNPAY   PaymentMethod = "PAYMENT_METHOD_VNPAY"
	PaymentMethodMOMO    PaymentMethod = "PAYMENT_METHOD_MOMO"
)

type Payment struct {
	ID          int64                   `json:"id"` /* unique */
	AccountID   int64                   `json:"account_id"`
	Method      PaymentMethod           `json:"method"`
	Status      PaymentStatus           `json:"status"`
	Total       commonmodel.Concurrency `json:"total"`
	DateCreated time.Time               `json:"date_created"`
}

type PaymentItem struct {
	ID        int64                   `json:"id"` /* unique */
	PaymentID int64                   `json:"payment_id"`
	Name      string                  `json:"name"`
	Price     commonmodel.Concurrency `json:"price"`
}

type PaymentVNPAY struct {
	ID                 int64  `json:"id"` /* unique */
	VnpTxnRef          string `json:"vnp_txn_ref"`
	VnpOrderInfo       string `json:"vnp_order_info"`
	VnpTransactionNo   string `json:"vnp_transaction_no"`
	VnpTransactionDate string `json:"vnp_transaction_date"`
	VnpCreateDate      string `json:"vnp_create_date"`
	VnpIpAddr          string `json:"vnp_ip_addr"`
}

type PaymentProcesseDataNATS struct {
	PaymentID int64 `json:"paymentID"`
}
