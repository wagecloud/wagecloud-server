package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	paymentservice "github.com/wagecloud/wagecloud-server/internal/modules/payment/service"
	paymentstorage "github.com/wagecloud/wagecloud-server/internal/modules/payment/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type EchoHandler struct {
	service paymentservice.Service
}

func NewEchoHandler(service paymentservice.Service) *EchoHandler {
	return &EchoHandler{service: service}
}

type GetPaymentRequest struct {
	ID int64 `param:"id" validate:"required"`
}

func (h *EchoHandler) GetPayment(c echo.Context) error {
	var req GetPaymentRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	payment, err := h.service.GetPayment(c.Request().Context(), req.ID)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, payment)
}

type ListPaymentsRequest struct {
	Page            int32                       `query:"page" validate:"min=1"`
	Limit           int32                       `query:"limit" validate:"min=5,max=100"`
	ID              *string                     `query:"id"`
	AccountID       *int64                      `query:"account_id"`
	Method          *paymentmodel.PaymentMethod `query:"method"`
	Status          *paymentmodel.PaymentStatus `query:"status"`
	DateCreatedFrom *int64                      `query:"date_created_from"`
	DateCreatedTo   *int64                      `query:"date_created_to"`
}

func (h *EchoHandler) ListPayments(c echo.Context) error {
	var req ListPaymentsRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	params := paymentstorage.ListPaymentsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		ID:              req.ID,
		AccountID:       req.AccountID,
		Method:          req.Method,
		Status:          req.Status,
		DateCreatedFrom: req.DateCreatedFrom,
		DateCreatedTo:   req.DateCreatedTo,
	}

	result, err := h.service.ListPayments(c.Request().Context(), params)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, result)
}

type CreatePaymentRequest struct {
}

func (h *EchoHandler) CreatePayment(c echo.Context) error {
	var req CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	// TODO: fix the payment create service
	// payment, err := h.service.CreatePayment(c.Request().Context(), )
	// if err != nil {
	// 	return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	// }

	// return response.FromDTO(c.Response().Writer, http.StatusCreated, payment)
	return nil
}

type UpdatePaymentRequest struct {
	ID     int64                       `param:"id" validate:"required"`
	Method *paymentmodel.PaymentMethod `json:"method"`
	Status *paymentmodel.PaymentStatus `json:"status"`
	Total  *int64                      `json:"total"`
}

func (h *EchoHandler) UpdatePayment(c echo.Context) error {
	var req UpdatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	params := paymentstorage.UpdatePaymentParams{
		ID:     req.ID,
		Method: req.Method,
		Status: req.Status,
		Total:  req.Total,
	}

	payment, err := h.service.UpdatePayment(c.Request().Context(), params)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, payment)
}

type DeletePaymentRequest struct {
	ID int64 `param:"id" validate:"required"`
}

func (h *EchoHandler) DeletePayment(c echo.Context) error {
	var req DeletePaymentRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := h.service.DeletePayment(c.Request().Context(), req.ID); err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Payment deleted successfully")
}

type CreatePaymentItemRequest struct {
	PaymentID int64  `param:"id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Price     int64  `json:"price" validate:"required"`
}

func (h *EchoHandler) CreatePaymentItem(c echo.Context) error {
	var req CreatePaymentItemRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromDTO(c.Response().Writer, http.StatusBadRequest, err)
	}

	item := paymentmodel.PaymentItem{
		PaymentID: req.PaymentID,
		Name:      req.Name,
		Price:     req.Price,
	}

	result, err := h.service.CreatePaymentItem(c.Request().Context(), item)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, result)
}

type CreatePaymentVNPAYRequest struct {
	ID                 int64  `param:"id" validate:"required"`
	VnpTxnRef          string `json:"vnp_txn_ref" validate:"required"`
	VnpOrderInfo       string `json:"vnp_order_info" validate:"required"`
	VnpTransactionNo   string `json:"vnp_transaction_no" validate:"required"`
	VnpTransactionDate string `json:"vnp_transaction_date" validate:"required"`
	VnpCreateDate      string `json:"vnp_create_date" validate:"required"`
	VnpIpAddr          string `json:"vnp_ip_addr" validate:"required"`
}

func (h *EchoHandler) CreatePaymentVNPAY(c echo.Context) error {
	var req CreatePaymentVNPAYRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	vnpay := paymentmodel.PaymentVNPAY{
		ID:                 req.ID,
		VnpTxnRef:          req.VnpTxnRef,
		VnpOrderInfo:       req.VnpOrderInfo,
		VnpTransactionNo:   req.VnpTransactionNo,
		VnpTransactionDate: req.VnpTransactionDate,
		VnpCreateDate:      req.VnpCreateDate,
		VnpIpAddr:          req.VnpIpAddr,
	}

	result, err := h.service.CreatePaymentVNPAY(c.Request().Context(), vnpay)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, result)
}
