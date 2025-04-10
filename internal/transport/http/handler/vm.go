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
		Role:      claims.Role,
		AccountID: claims.AccountID,
		ID:        vmID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, vm, http.StatusOK)
}

type ListVMsRequest struct {
	Page          int32   `schema:"page" validate:"required,min=1"`
	Limit         int32   `schema:"limit" validate:"required,min=5,max=100"`
	NetworkID     *string `schema:"network_id" validate:"omitempty"`
	OsID          *string `schema:"os_id" validate:"omitempty"`
	ArchID        *string `schema:"arch_id" validate:"omitempty"`
	Name          *string `schema:"name" validate:"omitempty,min=1,max=255"`
	CpuFrom       *int64  `schema:"cpu_from" validate:"omitempty,min=1,ltefield=CpuTo"`
	CpuTo         *int64  `schema:"cpu_to" validate:"omitempty,min=1,gtefield=CpuFrom"`
	RamFrom       *int64  `schema:"ram_from" validate:"omitempty,min=1,ltefield=RamTo"`
	RamTo         *int64  `schema:"ram_to" validate:"omitempty,min=1,gtefield=RamFrom"`
	StorageFrom   *int64  `schema:"storage_from" validate:"omitempty,min=1,ltefield=StorageTo"`
	StorageTo     *int64  `schema:"storage_to" validate:"omitempty,min=1,gtefield=StorageFrom"`
	CreatedAtFrom *int64  `schema:"created_at_from" validate:"omitempty,min=0,ltefield=CreatedAtTo"`
	CreatedAtTo   *int64  `schema:"created_at_to" validate:"omitempty,min=0,gtefield=CreatedAtFrom"`
}

func (h *Handler) ListVMs(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.FromHTTPError(w, http.StatusUnauthorized)
		return
	}

	var req ListVMsRequest

	if err := decodeAndValidate(&req, r.URL.Query()); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	vms, err := h.service.VM.ListVMs(r.Context(), vm.ListVMsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		AccountID:     claims.AccountID,
		Role:          claims.Role,
		NetworkID:     req.NetworkID,
		OsID:          req.OsID,
		ArchID:        req.ArchID,
		Name:          req.Name,
		CpuFrom:       req.CpuFrom,
		CpuTo:         req.CpuTo,
		RamFrom:       req.RamFrom,
		RamTo:         req.RamTo,
		StorageFrom:   req.StorageFrom,
		StorageTo:     req.StorageTo,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, vms, http.StatusOK)
}

type CreateVMRequest struct {
	Userdata struct {
		Name              string   `json:"name" schema:"name" validate:"required,min=1,max=255"`
		SSHAuthorizedKeys []string `json:"ssh-authorized-keys" schema:"ssh-authorized-keys" validate:"omitempty,dive,min=20,max=5000"`
		Password          string   `json:"password" schema:"password" validate:"required,min=8,max=72"`
	} `json:"userdata" schema:"userdata" validate:"required"`
	Metadata struct {
		LocalHostname string `json:"local-hostname" schema:"local-hostname" validate:"required,hostname"`
	} `json:"metadata" schema:"metadata" validate:"required"`
	Spec struct {
		OsID    string `json:"os_id" schema:"os_id" validate:"required,min=1,max=255"`
		ArchID  string `json:"arch_id" schema:"arch_id" validate:"required,min=1,max=255"`
		Memory  int    `json:"memory" schema:"memory" validate:"required,min=512,max=262144"` // 512MB to 262144MB (256GB)
		Cpu     int    `json:"cpu" schema:"cpu" validate:"required,min=1,max=64"`             // 1 to 64 cores
		Storage int    `json:"storage" schema:"storage" validate:"required,min=10,max=2048"`  // 10GB to 2TB
	} `json:"spec" schema:"spec" validate:"required"`
}

func (h *Handler) CreateVM(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetClaims(r)
	if err != nil {
		response.FromHTTPError(w, http.StatusUnauthorized)
		return
	}

	var req CreateVMRequest
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	vm, err := h.service.VM.CreateVM(r.Context(), vm.CreateVMParams{
		AccountID:         claims.AccountID,
		Name:              req.Userdata.Name,
		SSHAuthorizedKeys: req.Userdata.SSHAuthorizedKeys,
		Password:          req.Userdata.Password,
		LocalHostname:     req.Metadata.LocalHostname,
		OsID:              req.Spec.OsID,
		ArchID:            req.Spec.ArchID,
		Memory:            req.Spec.Memory,
		Cpu:               req.Spec.Cpu,
		Storage:           req.Spec.Storage,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, vm, http.StatusCreated)
}

type UpdateVMRequest struct {
	NetworkID *string `json:"network_id" schema:"network_id" validate:"omitempty"`
	OsID      *string `json:"os_id" schema:"os_id" validate:"omitempty"`
	ArchID    *string `json:"arch_id" schema:"arch_id" validate:"omitempty"`
	Name      *string `json:"name" schema:"name" validate:"omitempty,min=1,max=255"`
	Cpu       *int32  `json:"cpu" schema:"cpu" validate:"omitempty,min=1,max=64"`
	Ram       *int32  `json:"ram" schema:"ram" validate:"omitempty,min=512,max=262144"`
	Storage   *int32  `json:"storage" schema:"storage" validate:"omitempty,min=10,max=2048"`
}

func (h *Handler) UpdateVM(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateVMRequest
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	vm, err := h.service.VM.UpdateVM(r.Context(), vm.UpdateVMParams{
		Role:      claims.Role,
		AccountID: claims.AccountID,
		ID:        vmID,
		NetworkID: req.NetworkID,
		OsID:      req.OsID,
		ArchID:    req.ArchID,
		Name:      req.Name,
		Cpu:       req.Cpu,
		Ram:       req.Ram,
		Storage:   req.Storage,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, vm, http.StatusOK)
}

func (h *Handler) DeleteVM(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.VM.DeleteVM(r.Context(), vm.DeleteVMParams{
		Role:      claims.Role,
		AccountID: claims.AccountID,
		ID:        vmID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, nil, http.StatusOK)
}

func (h *Handler) StartVM(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.VM.StartVM(r.Context(), vm.StartVMParams{
		AccountID: claims.AccountID,
		Role:      claims.Role,
		ID:        vmID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, nil, http.StatusOK)
}

func (h *Handler) StopVM(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.VM.StopVM(r.Context(), vm.StopVMParams{
		AccountID: claims.AccountID,
		Role:      claims.Role,
		ID:        vmID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, nil, http.StatusOK)
}
