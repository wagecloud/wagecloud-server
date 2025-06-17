package accountecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type EchoHandler struct {
	service accountsvc.Service
}

func NewEchoHandler(service accountsvc.Service) *EchoHandler {
	return &EchoHandler{service: service}
}

type GetUserRequest struct {
	ID       *int64  `query:"id" validate:"omitempty"`
	Username *string `query:"username" validate:"omitempty,min=1,max=255"`
	Email    *string `query:"email" validate:"omitempty,email"`
	Phone    *string `query:"phone" validate:"omitempty,phone"`
}

func (h *EchoHandler) GetUser(c echo.Context) error {
	var req GetUserRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusUnauthorized, err)
	}

	if req.ID == nil && req.Username == nil && req.Email == nil {
		req.ID = &claims.AccountID
	}

	account, err := h.service.GetUser(c.Request().Context(), accountsvc.GetUserParams{
		Account:  claims.ToAuthenticatedAccount(),
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, account)
}

type LoginUserRequest struct {
	ID       *int64  `json:"id" validate:"omitempty"`
	Username *string `json:"username" validate:"omitempty"`
	Email    *string `json:"email" validate:"omitempty"`
	Phone    *string `json:"phone" validate:"omitempty"`
	Password string  `json:"password" validate:"required"`
}

func (h *EchoHandler) LoginUser(c echo.Context) error {
	var req LoginUserRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	account, err := h.service.LoginUser(c.Request().Context(), accountsvc.LoginUserParams{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusUnauthorized, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, account)
}

type RegisterUserRequest struct {
	FirstName string  `json:"first_name" validate:"required,min=1,max=255"`
	LastName  string  `json:"last_name" validate:"required,min=1,max=255"`
	Username  string  `json:"username" validate:"required,min=1,max=255"`
	Password  string  `json:"password" validate:"required,min=8,max=72"`
	Email     *string `json:"email" validate:"required,email"`
	Phone     *string `json:"phone" validate:"omitempty,phone"`
}

func (h *EchoHandler) RegisterUser(c echo.Context) error {
	var req RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	result, err := h.service.RegisterUser(c.Request().Context(), accountsvc.RegisterUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Phone:     req.Phone,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, result)
}
