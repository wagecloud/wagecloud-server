package handler

import (
	"encoding/json"
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
		response.FromError(w, err)
	}

	response.FromDTO(w, account, http.StatusOK, "Account retrieved successfully")
}

type LoginUserParams struct {
	ID       *int64  `json:"id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	account, err := h.service.Account.LoginUser(r.Context(), account.LoginUserParams{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, account, http.StatusOK, "Login successful")
}

type RegisterUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	result, err := h.service.Account.RegisterUser(r.Context(), model.AccountUser{
		AccountBase: model.AccountBase{
			Username: req.Username,
			Email:    req.Email,
			Name:     req.Name,
			Password: req.Password,
		},
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, result, http.StatusCreated, "User registered successfully")
}
