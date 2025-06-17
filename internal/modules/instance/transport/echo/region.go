package instanceecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type GetRegionRequest struct {
	ID string `param:"id" validate:"required"`
}

func (h *EchoHandler) GetRegion(c echo.Context) error {
	var req GetRegionRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	region, err := h.service.GetRegion(c.Request().Context(), req.ID)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, region)
}

type ListRegionsRequest struct {
	Page  int32   `query:"page" validate:"min=1"`
	Limit int32   `query:"limit" validate:"min=5,max=100"`
	ID    *string `query:"id"`
	Name  *string `query:"name"`
}

func (h *EchoHandler) ListRegions(c echo.Context) error {
	var req ListRegionsRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	regions, err := h.service.ListRegions(c.Request().Context(), instancesvc.ListRegionsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, regions)
}

type CreateRegionRequest struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (h *EchoHandler) CreateRegion(c echo.Context) error {
	var req CreateRegionRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	region, err := h.service.CreateRegion(c.Request().Context(), instancesvc.CreateRegionParams{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, region)
}

type UpdateRegionRequest struct {
	ID    string  `param:"id" validate:"required"`
	NewID *string `json:"new_id"`
	Name  *string `json:"name"`
}

func (h *EchoHandler) UpdateRegion(c echo.Context) error {
	var req UpdateRegionRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	region, err := h.service.UpdateRegion(c.Request().Context(), instancesvc.UpdateRegionParams{
		ID:    req.ID,
		NewID: req.NewID,
		Name:  req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, region)
}

type DeleteRegionRequest struct {
	ID string `param:"id" validate:"required"`
}

func (h *EchoHandler) DeleteRegion(c echo.Context) error {
	var req DeleteRegionRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := h.service.DeleteRegion(c.Request().Context(), req.ID); err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Region deleted successfully")
}
