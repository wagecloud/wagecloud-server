package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/network"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetNetwork(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"networkID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "networkID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	network, err := h.service.Network.GetNetwork(r.Context(), network.GetNetworkParams{
		ID: params.ID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, network, http.StatusOK)
}

func (h *Handler) ListNetworks(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Page          int32   `schema:"page" validate:"required,min=1"`
		Limit         int32   `schema:"limit" validate:"required,min=5,max=100"`
		ID            *string `schema:"id" validate:"omitempty"`
		PrivateIP     *string `schema:"private_ip" validate:"omitempty,ip"`
		CreatedAtFrom *int64  `schema:"created_at_from" validate:"omitempty,min=0,ltefield=CreatedAtTo"`
		CreatedAtTo   *int64  `schema:"created_at_to" validate:"omitempty,min=0,gtefield=CreatedAtFrom"`
	}
	if err := decodeAndValidate(&req, r.URL.Query()); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	networks, err := h.service.Network.ListNetworks(r.Context(), network.ListNetworksParams{
		PaginationParams: model.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		ID:            req.ID,
		PrivateIP:     req.PrivateIP,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromPaginate(w, networks)
}

func (h *Handler) CreateNetwork(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID        string `json:"id" validate:"required,min=1,max=255"`
		PrivateIP string `json:"private_ip" validate:"required,ip"`
	}
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	network, err := h.service.Network.CreateNetwork(r.Context(), network.CreateNetworkParams{
		ID:        req.ID,
		PrivateIP: req.PrivateIP,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, network, http.StatusCreated)
}

func (h *Handler) UpdateNetwork(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"networkID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "networkID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	var req struct {
		NewID     *string `json:"new_id" validate:"omitempty,min=1,max=255"`
		PrivateIP *string `json:"private_ip" validate:"omitempty,ip"`
	}
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	network, err := h.service.Network.UpdateNetwork(r.Context(), network.UpdateNetworkParams{
		ID:        params.ID,
		NewID:     req.NewID,
		PrivateIP: req.PrivateIP,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, network, http.StatusOK)
}

func (h *Handler) DeleteNetwork(w http.ResponseWriter, r *http.Request) {
	var params = struct {
		ID string `schema:"networkID" validate:"required,min=1,max=255"`
	}{
		ID: chi.URLParam(r, "networkID"),
	}
	if err := validate.Struct(params); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	err := h.service.Network.DeleteNetwork(r.Context(), network.DeleteNetworkParams{
		ID: params.ID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, nil, http.StatusOK)
}
