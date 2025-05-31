package osecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	osmodel "github.com/wagecloud/wagecloud-server/internal/modules/os/model"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

func (h *EchoHandler) GetArch(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	arch, err := h.osSvc.GetArch(c.Request().Context(), req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, arch)
}

func (h *EchoHandler) ListArchs(c echo.Context) error {
	var req struct {
		Page          int32   `query:"page" validate:"required,min=1"`
		Limit         int32   `query:"limit" validate:"required,min=5,max=100"`
		ID            *string `query:"id"`
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

	archs, err := h.osSvc.ListArchs(c.Request().Context(), ossvc.ListArchsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		ID:            req.ID,
		Name:          req.Name,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, archs)
}

func (h *EchoHandler) CreateArch(c echo.Context) error {
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

	arch, err := h.osSvc.CreateArch(c.Request().Context(), osmodel.Arch{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, arch)
}

func (h *EchoHandler) UpdateArch(c echo.Context) error {
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

	arch, err := h.osSvc.UpdateArch(c.Request().Context(), ossvc.UpdateArchParams{
		ID:    req.ID,
		NewID: req.NewID,
		Name:  req.Name,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, arch)
}

func (h *EchoHandler) DeleteArch(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.osSvc.DeleteArch(c.Request().Context(), req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
