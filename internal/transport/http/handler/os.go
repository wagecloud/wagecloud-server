package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/os"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetOS(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"osID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "osID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	os, err := h.service.OS.GetOS(r.Context(), os.GetOSParams{
		ID: params.ID,
	})

	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, os, http.StatusOK)
}

func (h *Handler) ListOSs(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Page          int32   `schema:"page" validate:"required,min=1"`
		Limit         int32   `schema:"limit" validate:"required,min=5,max=100"`
		CreatedAtFrom *int64  `schema:"created_at_from" validate:"omitempty,min=0,ltefield=CreatedAtTo"`
		CreatedAtTo   *int64  `schema:"created_at_to" validate:"omitempty,min=0,gtefield=CreatedAtFrom"`
		Name          *string `schema:"name" validate:"omitempty,min=1,max=255"`
	}
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

func (h *Handler) CreateOS(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   string `json:"id" validate:"required,min=1,max=255"`
		Name string `json:"name" validate:"required,min=1,max=255"`
	}
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

func (h *Handler) UpdateOS(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"osID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "osID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	var req struct {
		NewID *string `json:"new_id" validate:"omitempty,min=1,max=255"`
		Name  *string `json:"name" validate:"omitempty,min=1,max=255"`
	}
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	os, err := h.service.OS.UpdateOS(r.Context(), os.UpdateOSParams{
		ID:    params.ID,
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
	var params = struct {
		ID string `schema:"osID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "osID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.OS.DeleteOS(r.Context(), os.DeleteOSParams{
		ID: params.ID,
	}); err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, nil, http.StatusOK)
}
