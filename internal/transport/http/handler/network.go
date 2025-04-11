package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/service/network"
	"github.com/wagecloud/wagecloud-server/internal/transport/http/response"
)

func (h *Handler) GetNetwork(w http.ResponseWriter, r *http.Request) {
	networkID := chi.URLParam(r, "networkID")
	if networkID == "" {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	network, err := h.service.Network.GetNetwork(r.Context(), network.GetNetworkParams{
		ID: networkID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, network, http.StatusOK)
}

type ListNetworksRequest struct {
	Page          int32   `validate:"required,min=1"`
	Limit         int32   `validate:"required,min=5,max=100"`
	ID            *string `validate:"omitempty"`
	PrivateIP     *string `validate:"omitempty,ip"`
	CreatedAtFrom *int64  `validate:"omitempty,min=0,ltefield=CreatedAtTo"`
	CreatedAtTo   *int64  `validate:"omitempty,min=0,gtefield=CreatedAtFrom"`
}

func (h *Handler) ListNetworks(w http.ResponseWriter, r *http.Request) {
	var req ListNetworksRequest
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

type CreateNetworkRequest struct {
	ID        string `json:"id" validate:"required,min=1,max=255"`
	PrivateIP string `json:"private_ip" validate:"required,ip"`
}

func (h *Handler) CreateNetwork(w http.ResponseWriter, r *http.Request) {
	var req CreateNetworkRequest
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

type UpdateNetworkRequest struct {
	NewID     *string `json:"new_id" validate:"omitempty,min=1,max=255"`
	PrivateIP *string `json:"private_ip" validate:"omitempty,ip"`
}

func (h *Handler) UpdateNetwork(w http.ResponseWriter, r *http.Request) {
	networkID := chi.URLParam(r, "networkID")
	if networkID == "" {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	var req UpdateNetworkRequest
	if err := decodeAndValidateJSON(&req, r.Body); err != nil {
		response.FromError(w, err, http.StatusBadRequest)
		return
	}

	network, err := h.service.Network.UpdateNetwork(r.Context(), network.UpdateNetworkParams{
		ID:        networkID,
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
	networkID := chi.URLParam(r, "networkID")
	if networkID == "" {
		response.FromHTTPError(w, http.StatusBadRequest)
		return
	}

	err := h.service.Network.DeleteNetwork(r.Context(), network.DeleteNetworkParams{
		ID: networkID,
	})
	if err != nil {
		response.FromError(w, err, http.StatusInternalServerError)
		return
	}

	response.FromDTO(w, nil, http.StatusOK)
}
