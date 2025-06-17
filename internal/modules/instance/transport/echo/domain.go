package instanceecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type GetDomainRequest struct {
	ID int64 `param:"id" validate:"required"`
}

func (h *EchoHandler) GetDomain(c echo.Context) error {
	var req GetDomainRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	domain, err := h.service.GetDomain(c.Request().Context(), req.ID)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, domain)
}

type ListDomainsRequest struct {
	Page      int32   `query:"page" validate:"min=1"`
	Limit     int32   `query:"limit" validate:"min=5,max=100"`
	NetworkID *int64  `query:"network_id"`
	Name      *string `query:"name"`
}

func (h *EchoHandler) ListDomains(c echo.Context) error {
	var req ListDomainsRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	domains, err := h.service.ListDomains(c.Request().Context(), instancesvc.ListDomainsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		NetworkID: req.NetworkID,
		Name:      req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, domains)
}

type CreateDomainRequest struct {
	NetworkID int64  `json:"network_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

func (h *EchoHandler) CreateDomain(c echo.Context) error {
	var req CreateDomainRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	domain, err := h.service.CreateDomain(c.Request().Context(), instancesvc.CreateDomainParams{
		NetworkID: req.NetworkID,
		Name:      req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, domain)
}

type UpdateDomainRequest struct {
	ID   int64   `param:"id" validate:"required"`
	Name *string `json:"name"`
}

func (h *EchoHandler) UpdateDomain(c echo.Context) error {
	var req UpdateDomainRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	domain, err := h.service.UpdateDomain(c.Request().Context(), instancesvc.UpdateDomainParams{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, domain)
}

type DeleteDomainRequest struct {
	ID int64 `param:"id" validate:"required"`
}

func (h *EchoHandler) DeleteDomain(c echo.Context) error {
	var req DeleteDomainRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := h.service.DeleteDomain(c.Request().Context(), req.ID); err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Domain deleted successfully")
}
