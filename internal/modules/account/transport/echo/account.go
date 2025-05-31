package accountecho

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
)

type EchoHandler struct {
	accountsvc accountsvc.Service
}

func NewEchoHandler(accountsvc accountsvc.Service) *EchoHandler {
	return &EchoHandler{
		accountsvc: accountsvc,
	}
}

func (h *EchoHandler) GetAccount(c echo.Context) error {
	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	account, err := h.accountsvc.GetAccount(c.Request().Context(), accountsvc.GetAccountParams{
		ID: &claims.AccountID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get account")
	}

	return c.JSON(http.StatusOK, account)
}

type LoginUserParams struct {
	ID       *int64  `json:"id" validate:"omitempty"`
	Username *string `json:"username" validate:"omitempty"`
	Email    *string `json:"email" validate:"omitempty"`
	Password string  `json:"password" validate:"required"`
}

func (h *EchoHandler) LoginUser(c echo.Context) error {
	var req LoginUserParams
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	account, err := h.accountsvc.LoginUser(c.Request().Context(), accountsvc.LoginUserParams{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	return c.JSON(http.StatusOK, account)
}

type RegisterUserParams struct {
	Username string `json:"username" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *EchoHandler) RegisterUser(c echo.Context) error {
	var req RegisterUserParams
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	result, err := h.accountsvc.RegisterUser(c.Request().Context(), accountmodel.AccountUser{
		AccountBase: accountmodel.AccountBase{
			Username: req.Username,
			Name:     req.Name,
			Password: req.Password,
		},
		Email: req.Email,
	})
	if err != nil {
		fmt.Println("Error registering user:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to register user")
	}

	return c.JSON(http.StatusCreated, result)
}
