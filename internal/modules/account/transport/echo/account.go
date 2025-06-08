package accountecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/http/response"
)

type EchoHandler struct {
	service accountsvc.Service
}

func NewEchoHandler(service accountsvc.Service) *EchoHandler {
	return &EchoHandler{service: service}
}

type GetAccountRequest struct {
	ID *int64 `param:"id" validate:"omitempty"`
}

func (h *EchoHandler) GetAccount(c echo.Context) error {
	var req GetAccountRequest
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

	account, err := h.service.GetAccount(c.Request().Context(), accountsvc.GetAccountParams{
		ID: &claims.AccountID,
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
		Password: req.Password,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusUnauthorized, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, account)
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *EchoHandler) RegisterUser(c echo.Context) error {
	var req RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	result, err := h.service.RegisterUser(c.Request().Context(), accountmodel.AccountUser{
		AccountBase: accountmodel.AccountBase{
			Username: req.Username,
			Name:     req.Name,
			Password: req.Password,
		},
		Email: req.Email,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, result)
}
