package instanceecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type GetNetworkRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) GetNetwork(c echo.Context) error {
	var req GetNetworkRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	network, err := h.service.GetNetwork(c.Request().Context(), instancesvc.GetNetworkParams{
		ID: req.ID,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, network)
}

type ListNetworksRequest struct {
	Page          int32   `query:"page" validate:"min=1"`
	Limit         int32   `query:"limit" validate:"min=5,max=100"`
	ID            *string `query:"id"`
	PrivateIP     *string `query:"private_ip"`
	CreatedAtFrom *int64  `query:"created_at_from"`
	CreatedAtTo   *int64  `query:"created_at_to"`
}

func (h *EchoHandler) ListNetworks(c echo.Context) error {
	var req ListNetworksRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	networks, err := h.service.ListNetworks(c.Request().Context(), instancesvc.ListNetworksParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		ID:            req.ID,
		PrivateIP:     req.PrivateIP,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, networks)
}

type CreateNetworkRequest struct {
	ID        string `json:"id" validate:"required,min=1,max=255"`
	PrivateIP string `json:"private_ip" validate:"required,ip"`
}

func (h *EchoHandler) CreateNetwork(c echo.Context) error {
	var req CreateNetworkRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	network, err := h.service.CreateNetwork(c.Request().Context(), instancesvc.CreateNetworkParams{
		ID:        req.ID,
		PrivateIP: req.PrivateIP,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, network)
}

type UpdateNetworkRequest struct {
	ID        string  `param:"id" validate:"required,min=1,max=255"`
	NewID     *string `json:"new_id"`
	PrivateIP *string `json:"private_ip"`
}

func (h *EchoHandler) UpdateNetwork(c echo.Context) error {
	var req UpdateNetworkRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	network, err := h.service.UpdateNetwork(c.Request().Context(), instancesvc.UpdateNetworkParams{
		ID:        req.ID,
		PrivateIP: req.PrivateIP,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, network)
}

type DeleteNetworkRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) DeleteNetwork(c echo.Context) error {
	var req DeleteNetworkRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	err := h.service.DeleteNetwork(c.Request().Context(), instancesvc.DeleteNetworkParams{
		ID: req.ID,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Network deleted successfully")
}

