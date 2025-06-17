package osecho

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type GetArchRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) GetArch(c echo.Context) error {
	var req GetArchRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	arch, err := h.service.GetArch(c.Request().Context(), req.ID)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, arch)
}

type ListArchsRequest struct {
	Page          int32   `query:"page" validate:"min=1"`
	Limit         int32   `query:"limit" validate:"min=5,max=100"`
	ID            *string `query:"id"`
	Name          *string `query:"name"`
	CreatedAtFrom *int64  `query:"created_at_from"`
	CreatedAtTo   *int64  `query:"created_at_to"`
}

func (h *EchoHandler) ListArchs(c echo.Context) error {
	var req ListArchsRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	archs, err := h.service.ListArchs(c.Request().Context(), ossvc.ListArchsParams{
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
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	js, _ := json.Marshal(archs)
	fmt.Println("ListArchs response:", string(js))

	return response.FromPaginate(c.Response().Writer, archs)
}

type CreateArchRequest struct {
	ID   string `json:"id" validate:"required,min=1,max=255"`
	Name string `json:"name" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) CreateArch(c echo.Context) error {
	var req CreateArchRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	arch, err := h.service.CreateArch(c.Request().Context(), ossvc.CreateArchParams{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, arch)
}

type UpdateArchRequest struct {
	ID    string  `param:"id" validate:"required,min=1,max=255"`
	NewID *string `json:"new_id"`
	Name  *string `json:"name"`
}

func (h *EchoHandler) UpdateArch(c echo.Context) error {
	var req UpdateArchRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	arch, err := h.service.UpdateArch(c.Request().Context(), ossvc.UpdateArchParams{
		ID:    req.ID,
		NewID: req.NewID,
		Name:  req.Name,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, arch)
}

type DeleteArchRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) DeleteArch(c echo.Context) error {
	var req DeleteArchRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	err := h.service.DeleteArch(c.Request().Context(), req.ID)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Architecture deleted successfully")
}
