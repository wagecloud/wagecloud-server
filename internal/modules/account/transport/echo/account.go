package accountecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/echo/middleware/auth"
	"github.com/wagecloud/wagecloud-server/internal/shared/echo/response"
)

type EchoHandler struct {
	accountsvc accountsvc.Service
}

func (h *EchoHandler) GetAccount(c echo.Context) {
	claims, err := auth.GetClaims(c.Request())
	if err != nil {
		response.FromHTTPError(c.Response().Writer, http.StatusUnauthorized)
		return
	}

	account, err := h.accountsvc.GetAccount(c.Request().Context(), accountsvc.GetAccountParams{
		ID: &claims.AccountID,
	})
	if err != nil {
		response.FromError(c.Response().Writer, err, http.StatusInternalServerError)
	}

	response.FromDTO(c.Response().Writer, account, http.StatusOK)
}

type LoginUserParams struct {
	ID       *int64  `json:"id" validate:"omitempty"`
	Username *string `json:"username" validate:"omitempty"`
	Email    *string `json:"email" validate:"omitempty"`
	Password string  `json:"password" validate:"required"`
}

func (h *EchoHandler) LoginUser(c echo.Context) {
	var req LoginUserParams
	if err := c.Bind(&req); err != nil {
		response.FromError(c.Response().Writer, err, http.StatusBadRequest)
		return
	}

	account, err := h.accountsvc.LoginUser(c.Request().Context(), accountsvc.LoginUserParams{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.FromError(c.Response().Writer, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(c.Response().Writer, account, http.StatusOK)
}

type RegisterUserParams struct {
	Username string `json:"username" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *EchoHandler) RegisterUser(c echo.Context) {
	var req RegisterUserParams
	if err := c.Bind(&req); err != nil {
		response.FromError(c.Response().Writer, err, http.StatusBadRequest)
		return
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
		response.FromError(c.Response().Writer, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(c.Response().Writer, result, http.StatusCreated)
}
