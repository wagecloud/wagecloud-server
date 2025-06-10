package osecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type EchoHandler struct {
	service ossvc.Service
}

func NewEchoHandler(service ossvc.Service) *EchoHandler {
	return &EchoHandler{service: service}
}

type GetOSRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) GetOS(c echo.Context) error {
	var req GetOSRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	os, err := h.service.GetOS(c.Request().Context(), ossvc.GetOSParams{
		ID: req.ID,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, os)
}

type ListOSsRequest struct {
	Page          int32   `query:"page" validate:"min=1"`
	Limit         int32   `query:"limit" validate:"min=5,max=100"`
	Name          *string `query:"name"`
	CreatedAtFrom *int64  `query:"created_at_from"`
	CreatedAtTo   *int64  `query:"created_at_to"`
}

func (h *EchoHandler) ListOSs(c echo.Context) error {
	var req ListOSsRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	osList, err := h.service.ListOSs(c.Request().Context(), ossvc.ListOSsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Name:          req.Name,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, osList)
}

type CreateOSRequest struct {
	ID   string `json:"id" validate:"required,min=1,max=255"`
	Name string `json:"name" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) CreateOS(c echo.Context) error {
	var req CreateOSRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	os, err := h.service.CreateOS(c.Request().Context(), ossvc.CreateOSParams{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, os)
}

type UpdateOSRequest struct {
	ID    string  `param:"id" validate:"required,min=1,max=255"`
	NewID *string `json:"new_id"`
	Name  *string `json:"name"`
}

func (h *EchoHandler) UpdateOS(c echo.Context) error {
	var req UpdateOSRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	os, err := h.service.UpdateOS(c.Request().Context(), ossvc.UpdateOSParams{
		ID:    req.ID,
		NewID: req.NewID,
		Name:  req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, os)
}

type DeleteOSRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) DeleteOS(c echo.Context) error {
	var req DeleteOSRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	err := h.service.DeleteOS(c.Request().Context(), ossvc.DeleteOSParams{
		ID: req.ID,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "OS deleted successfully")
}
