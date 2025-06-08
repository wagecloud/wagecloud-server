package paymentmodel

type PaymentMethod string

const (
	PaymentMethodVNPAY PaymentMethod = "VNPAY"
	PaymentMethodMOMO  PaymentMethod = "MOMO"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusSuccess  PaymentStatus = "SUCCESS"
	PaymentStatusCanceled PaymentStatus = "CANCELED"
	PaymentStatusFailed   PaymentStatus = "FAILED"
)

type PaymentBase struct {
	ID          int64         `json:"id"` /* unique */
	AccountID   int64         `json:"account_id"`
	Method      PaymentMethod `json:"method"`
	Status      PaymentStatus `json:"status"`
	Total       int64         `json:"total"`
	DateCreated int64         `json:"date_created"`
}

type PaymentItem struct {
	ID        int64  `json:"id"` /* unique */
	PaymentID int64  `json:"payment_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
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
