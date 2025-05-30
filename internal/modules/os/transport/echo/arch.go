package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/arch"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetArch(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"archID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "archID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	arch, err := h.service.Arch.GetArch(r.Context(), params.ID)
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, arch, http.StatusOK)
}

func (h *Handler) ListArchs(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Page          int32   `schema:"page" validate:"required,min=1"`
		Limit         int32   `schema:"limit" validate:"required,min=5,max=100"`
		ID            *string `schema:"id" validate:"omitempty"`
		Name          *string `schema:"name" validate:"omitempty,min=1,max=255"`
		CreatedAtFrom *int64  `schema:"created_at_from" validate:"omitempty,min=0,ltefield=CreatedAtTo"`
		CreatedAtTo   *int64  `schema:"created_at_to" validate:"omitempty,min=0,gtefield=CreatedAtFrom"`
	}
	if err := decodeAndValidate(&req, r.URL.Query()); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	archs, err := h.service.Arch.ListArchs(r.Context(), arch.ListArchsParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		ID:            req.ID,
		Name:          req.Name,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromPaginate(w, archs)
}

func (h *Handler) CreateArch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   string `json:"id" validate:"required,min=1,max=255"`
		Name string `json:"name" validate:"required,min=1,max=255"`
	}
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	arch, err := h.service.Arch.CreateArch(r.Context(), model.Arch{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, arch, http.StatusCreated)
}

func (h *Handler) UpdateArch(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"archID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "archID"),
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

	arch, err := h.service.Arch.UpdateArch(r.Context(), arch.UpdateArchParams{
		ID:    params.ID,
		NewID: req.NewID,
		Name:  req.Name,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, arch, http.StatusOK)
}

func (h *Handler) DeleteArch(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"archID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "archID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.Arch.DeleteArch(r.Context(), params.ID); err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, nil, http.StatusOK)
}
