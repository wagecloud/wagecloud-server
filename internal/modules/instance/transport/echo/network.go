package instanceecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

func (h *EchoHandler) GetNetwork(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	network, err := h.instanceSvc.GetNetwork(c.Request().Context(), instancesvc.GetNetworkParams{
		ID: req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, network)
}

func (h *EchoHandler) ListNetworks(c echo.Context) error {
	var req struct {
		Page          int32   `query:"page" validate:"required,min=1"`
		Limit         int32   `query:"limit" validate:"required,min=5,max=100"`
		ID            *string `query:"id"`
		PrivateIP     *string `query:"private_ip"`
		CreatedAtFrom *int64  `query:"created_at_from"`
		CreatedAtTo   *int64  `query:"created_at_to"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	networks, err := h.instanceSvc.ListNetworks(c.Request().Context(), instancesvc.ListNetworksParams{
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, networks)
}

func (h *EchoHandler) CreateNetwork(c echo.Context) error {
	var req struct {
		ID        string `json:"id" validate:"required,min=1,max=255"`
		PrivateIP string `json:"private_ip" validate:"required,ip"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	network, err := h.instanceSvc.CreateNetwork(c.Request().Context(), instancesvc.CreateNetworkParams{
		ID:        req.ID,
		PrivateIP: req.PrivateIP,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, network)
}

func (h *EchoHandler) UpdateNetwork(c echo.Context) error {
	var req struct {
		ID        string  `param:"id" validate:"required,min=1,max=255"`
		NewID     *string `json:"new_id"`
		PrivateIP *string `json:"private_ip"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	network, err := h.instanceSvc.UpdateNetwork(c.Request().Context(), instancesvc.UpdateNetworkParams{
		ID:        req.ID,
		NewID:     req.NewID,
		PrivateIP: req.PrivateIP,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, network)
}

func (h *EchoHandler) DeleteNetwork(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.instanceSvc.DeleteNetwork(c.Request().Context(), instancesvc.DeleteNetworkParams{
		ID: req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
