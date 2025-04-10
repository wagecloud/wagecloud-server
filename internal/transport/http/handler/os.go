package handler

import (
	"encoding/json"
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
	}

	os, err := h.service.OS.GetOS(r.Context(), os.GetOSParams{})

	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, os, http.StatusOK)
}

type ListOSsRequest struct {
	Page          int32   `schema:"page" validate:"required,min=1"`
	Limit         int32   `schema:"limit" validate:"required,min=5,max=100"`
	CreatedAtFrom *int64  `schema:"created_at_from" validate:"omitempty,min=0,ltefield=CreatedAtTo"`
	CreatedAtTo   *int64  `schema:"created_at_to" validate:"omitempty,min=0,gtefield=CreatedAtFrom"`
	Name          *string `schema:"name" validate:"omitempty,min=1,max=255"`
}

func (h *Handler) ListOSs(w http.ResponseWriter, r *http.Request) {
	// claims, err := auth.GetClaims(r)

	// if err != nil {
	//   response.FromHTTPError(w, http.StatusUnauthorized)
	//   return
	// }

	var req ListOSsRequest

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
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

	response.FromDTO(w, osList, http.StatusOK)
}

type CreateOSRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

func (h *Handler) CreateOS(w http.ResponseWriter, r *http.Request) {
	var req CreateOSRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	os, err := h.service.OS.CreateOS(r.Context(), os.CreateOSParams{
		Name: req.Name,
	})

	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

  response.FromDTO(w, os, http.StatusCreated)
}

func (h *Handler) UpdateOS(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteOS(w http.ResponseWriter, r *http.Request) {
}
