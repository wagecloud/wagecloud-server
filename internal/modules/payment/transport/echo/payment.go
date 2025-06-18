package paymentecho

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wagecloud/wagecloud-server/internal/logger"
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

func (h *EchoHandler) RegisterRoutes(g *echo.Group) {
	fmt.Println("Registering payment routes")
	payment := g.Group("/payment")
	// handle vnpay ipn
	payment.GET("/vnpay/", h.VnpayVerifyIPN)

	payment.GET("/", h.ListPayments)
	payment.GET("/:id", h.GetPayment)
	payment.POST("/", h.CreatePayment)
	payment.PATCH("/:id", h.UpdatePayment)
	payment.DELETE("/:id", h.DeletePayment)

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

func (h *EchoHandler) VnpayVerifyIPN(c echo.Context) error {
	var query map[string]any
	fmt.Println("Received VNPAY IPN request")

	if err := c.Bind(&query); err != nil {
		logger.Log.Error("Failed to bind query parameters" + "error" + err.Error())
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	fmt.Println("Received VNPAY IPN:", query)

	// Verify the checksum hash
	if _, err := h.service.VerifyPayment(c.Request().Context(), paymentmodel.PaymentMethodVNPAY, query); err != nil {
		c.NoContent(http.StatusBadRequest)
		fmt.Println("Payment verification failed:", err)
		return nil
	}

	return c.NoContent(http.StatusOK)
}
