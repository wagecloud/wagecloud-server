package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/vm"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/middleware/auth"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetVM(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.FromHTTPError(w, http.StatusUnauthorized)
		return
	}

	vmID := chi.URLParam(r, "vmID")
	if vmID == "" {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	vm, err := h.service.VM.GetVM(r.Context(), vm.GetVMParams{
		AccountID: claims.AccountID,
		// ID:        vmID,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, vm, http.StatusOK, "VM retrieved successfully")
}

func (h *Handler) ListVMs(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.FromHTTPError(w, http.StatusUnauthorized)
		return
	}

	vms, err := h.service.VM.ListVMs(r.Context(), vm.ListVMsParams{
		PaginationParams: model.PaginationParams{
			Page:  1,
			Limit: 10,
		},
		AccountID: claims.AccountID,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}

	response.FromDTO(w, vms, http.StatusOK, "VMs retrieved successfully")
}
