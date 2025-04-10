package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/os"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetOS(w http.ResponseWriter, r *http.Request) {
	osID := chi.URLParam(r, "osID")

	if osID == "" {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	os, err := h.service.OS.GetOS(r.Context(), os.GetOSParams{
		ID: osID,
	})

	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, os, http.StatusOK)
}

type ListOSsRequest struct {
	Page          int32   `validate:"required,min=1"`
	Limit         int32   `validate:"required,min=5,max=100"`
	CreatedAtFrom *int64  `validate:"omitempty,min=0,ltefield=CreatedAtTo"`
	CreatedAtTo   *int64  `validate:"omitempty,min=0,gtefield=CreatedAtFrom"`
	Name          *string `validate:"omitempty,min=1,max=255"`
}

func (h *Handler) ListOSs(w http.ResponseWriter, r *http.Request) {
	var req ListOSsRequest

	if err := decodeAndValidate(&req, r.URL.Query()); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	osList, err := h.service.OS.ListOSs(r.Context(), os.ListOSsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Name:          req.Name,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})

	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromPaginate(w, osList)
}

type CreateOSRequest struct {
	ID   string `json:"id" validate:"required,min=1,max=255"`
	Name string `json:"name" validate:"required,min=1,max=255"`
}

func (h *Handler) CreateOS(w http.ResponseWriter, r *http.Request) {
	var req CreateOSRequest

	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	os, err := h.service.OS.CreateOS(r.Context(), os.CreateOSParams{
		ID:   req.ID,
		Name: req.Name,
	})

	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, os, http.StatusCreated)
}

type UpdateOSRequest struct {
	ID    string  `json:"id" validate:"required,min=1,max=255"`
	NewID *string `json:"new_id" validate:"omitempty,min=1,max=255"`
	Name  *string `json:"name" validate:"omitempty,min=1,max=255"`
}

func (h *Handler) UpdateOS(w http.ResponseWriter, r *http.Request) {
	var req UpdateOSRequest

	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	os, err := h.service.OS.UpdateOS(r.Context(), os.UpdateOSParams{
		ID:    req.ID,
		NewID: req.NewID,
		Name:  req.Name,
	})

	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, os, http.StatusOK)
}

func (h *Handler) DeleteOS(w http.ResponseWriter, r *http.Request) {
	osID := chi.URLParam(r, "osID")

	if osID == "" {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	if err := h.service.OS.DeleteOS(r.Context(), os.DeleteOSParams{
		ID: osID,
	}); err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromHTTPError(w, http.StatusNoContent)
}
