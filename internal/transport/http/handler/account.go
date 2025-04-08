package handler

import (
	"net/http"

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
