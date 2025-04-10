package handler

import (
	"net/http"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/account"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/middleware/auth"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.FromHTTPError(w, http.StatusUnauthorized)
		return
	}

	account, err := h.service.Account.GetAccount(r.Context(), account.GetAccountParams{
		ID: &claims.AccountID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
	}

	response.FromDTO(w, account, http.StatusOK)
}

type LoginUserParams struct {
	ID       *int64  `json:"id" validate:"omitempty"`
	Username *string `json:"username" validate:"omitempty"`
	Email    *string `json:"email" validate:"omitempty"`
	Password string  `json:"password" validate:"required"`
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserParams
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	account, err := h.service.Account.LoginUser(r.Context(), account.LoginUserParams{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, account, http.StatusOK)
}

type RegisterUserParams struct {
	Username string `json:"username" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserParams
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	result, err := h.service.Account.RegisterUser(r.Context(), model.AccountUser{
		AccountBase: model.AccountBase{
			Username: req.Username,
			Name:     req.Name,
			Password: req.Password,
		},
		Email: req.Email,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, result, http.StatusCreated)
}
