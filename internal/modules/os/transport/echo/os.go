package osecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

type EchoHandler struct {
	osSvc ossvc.Service
}

func NewEchoHandler(osSvc ossvc.Service) *EchoHandler {
	return &EchoHandler{
		osSvc: osSvc,
	}
}

func (h *EchoHandler) GetOS(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	os, err := h.osSvc.GetOS(c.Request().Context(), ossvc.GetOSParams{
		ID: req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, os)
}

func (h *EchoHandler) ListOSs(c echo.Context) error {
	var req struct {
		Page          int32   `query:"page" validate:"required,min=1"`
		Limit         int32   `query:"limit" validate:"required,min=5,max=100"`
		Name          *string `query:"name"`
		CreatedAtFrom *int64  `query:"created_at_from"`
		CreatedAtTo   *int64  `query:"created_at_to"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	osList, err := h.osSvc.ListOSs(c.Request().Context(), ossvc.ListOSsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Name:          req.Name,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, osList)
}

func (h *EchoHandler) CreateOS(c echo.Context) error {
	var req struct {
		ID   string `json:"id" validate:"required,min=1,max=255"`
		Name string `json:"name" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	os, err := h.osSvc.CreateOS(c.Request().Context(), ossvc.CreateOSParams{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, os)
}

func (h *EchoHandler) UpdateOS(c echo.Context) error {
	var req struct {
		ID    string  `param:"id" validate:"required,min=1,max=255"`
		NewID *string `json:"new_id"`
		Name  *string `json:"name"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	os, err := h.osSvc.UpdateOS(c.Request().Context(), ossvc.UpdateOSParams{
		ID:    req.ID,
		NewID: req.NewID,
		Name:  req.Name,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, os)
}

func (h *EchoHandler) DeleteOS(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.osSvc.DeleteOS(c.Request().Context(), ossvc.DeleteOSParams{
		ID: req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
