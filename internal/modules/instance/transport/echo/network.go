package instanceecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type MapPortNginxRequest struct {
	VMIP         string `json:"vm_ip" validate:"required,ip"`
	ExternalPort int32  `json:"external_port" validate:"required,min=1,max=65535"`
	InternalPort int32  `json:"internal_port" validate:"required,min=1,max=65535"`
	Type         string `json:"type" validate:"required,oneof=stream http"`
}

func (h *EchoHandler) MapPortNginx(c echo.Context) error {
	var req MapPortNginxRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	err := h.service.MapPortNginx(c.Request().Context(), instancesvc.MapPortNginxParams{
		VMIP:         req.VMIP,
		ExternalPort: req.ExternalPort,
		InternalPort: req.InternalPort,
		Type:         req.Type,
	})

	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Port mapped successfully")
}

type UnmapPortNginxRequest struct {
	ExternalPort int32  `json:"external_port" validate:"required,min=1,max=65535"`
	Type         string `json:"type" validate:"required,oneof=stream http"`
}

func (h *EchoHandler) UnmapPortNginx(c echo.Context) error {
	var req UnmapPortNginxRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	err := h.service.UnmapPortNginx(c.Request().Context(), instancesvc.UnmapPortNginxParams{
		ExternalPort: req.ExternalPort,
		ProtocolType: req.Type,
	})

	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Port unmapped successfully")
}

type GetNetworkRequest struct {
	ID int64 `param:"id" validate:"required,min=1,max=255"`
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
	Page      int32   `query:"page" validate:"min=1"`
	Limit     int32   `query:"limit" validate:"min=5,max=100"`
	ID        *string `query:"id"`
	PrivateIP *string `query:"private_ip"`
	PublicIP  *string `query:"public_ip"`
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
		ID:        req.ID,
		PrivateIP: req.PrivateIP,
		PublicIP:  req.PublicIP,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, networks)
}

// type CreateNetworkRequest struct {
// 	ID        string `json:"id" validate:"required,min=1,max=255"`
// 	PrivateIP string `json:"private_ip" validate:"required,ip"`
// }

// func (h *EchoHandler) CreateNetwork(c echo.Context) error {
// 	var req CreateNetworkRequest
// 	if err := c.Bind(&req); err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
// 	}

// 	if err := c.Validate(&req); err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
// 	}

// 	network, err := h.service.CreateNetwork(c.Request().Context(), instancesvc.CreateNetworkParams{
// 		InstanceID: req.ID,
// 		PrivateIP: req.PrivateIP,
// 	})
// 	if err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
// 	}

// 	return response.FromDTO(c.Response().Writer, http.StatusCreated, network)
// }

// type UpdateNetworkRequest struct {
// 	ID        string  `param:"id" validate:"required,min=1,max=255"`
// 	NewID     *string `json:"new_id"`
// 	PrivateIP *string `json:"private_ip"`
// }

// func (h *EchoHandler) UpdateNetwork(c echo.Context) error {
// 	var req UpdateNetworkRequest
// 	if err := c.Bind(&req); err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
// 	}

// 	if err := c.Validate(&req); err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
// 	}

// 	network, err := h.service.UpdateNetwork(c.Request().Context(), instancesvc.UpdateNetworkParams{
// 		ID:        req.ID,
// 		PrivateIP: req.PrivateIP,
// 	})
// 	if err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
// 	}

// 	return response.FromDTO(c.Response().Writer, http.StatusOK, network)
// }

// type DeleteNetworkRequest struct {
// 	ID string `param:"id" validate:"required,min=1,max=255"`
// }

// func (h *EchoHandler) DeleteNetwork(c echo.Context) error {
// 	var req DeleteNetworkRequest
// 	if err := c.Bind(&req); err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
// 	}

// 	if err := c.Validate(&req); err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
// 	}

// 	err := h.service.DeleteNetwork(c.Request().Context(), instancesvc.DeleteNetworkParams{
// 		ID: req.ID,
// 	})
// 	if err != nil {
// 		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
// 	}

// 	return response.FromMessage(c.Response().Writer, http.StatusOK, "Network deleted successfully")
// }
