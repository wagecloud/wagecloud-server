package instanceecho

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type GetInstanceLogRequest struct {
	ID int64 `param:"id" validate:"required"`
}

func (h *EchoHandler) GetInstanceLog(c echo.Context) error {
	var req GetInstanceLogRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	log, err := h.service.GetInstanceLog(c.Request().Context(), req.ID)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, log)
}

type ListInstanceLogsRequest struct {
	Page          int32                  `query:"page" validate:"min=1"`
	Limit         int32                  `query:"limit" validate:"min=5,max=100"`
	InstanceID    *string                `query:"instance_id"`
	Type          *instancemodel.LogType `query:"type"`
	Title         *string                `query:"title"`
	Description   *string                `query:"description"`
	CreatedAtFrom *time.Time             `query:"created_at_from"`
	CreatedAtTo   *time.Time             `query:"created_at_to"`
}

func (h *EchoHandler) ListInstanceLogs(c echo.Context) error {
	var req ListInstanceLogsRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	logs, err := h.service.ListInstanceLogs(c.Request().Context(), instancesvc.ListInstanceLogsParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		InstanceID:    req.InstanceID,
		Type:          req.Type,
		Title:         req.Title,
		Description:   req.Description,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, logs)
}

type CreateInstanceLogRequest struct {
	InstanceID  string                `json:"instance_id" validate:"required"`
	Type        instancemodel.LogType `json:"type" validate:"required"`
	Title       string                `json:"title" validate:"required"`
	Description *string               `json:"description"`
}

func (h *EchoHandler) CreateInstanceLog(c echo.Context) error {
	var req CreateInstanceLogRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	log, err := h.service.CreateInstanceLog(c.Request().Context(), instancesvc.CreateInstanceLogParams{
		InstanceID:  req.InstanceID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, log)
}

type UpdateInstanceLogRequest struct {
	ID              int64                  `param:"id" validate:"required"`
	Type            *instancemodel.LogType `json:"type"`
	Title           *string                `json:"title"`
	Description     *string                `json:"description"`
	NullDescription bool                   `json:"null_description"`
}

func (h *EchoHandler) UpdateInstanceLog(c echo.Context) error {
	var req UpdateInstanceLogRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	log, err := h.service.UpdateInstanceLog(c.Request().Context(), instancesvc.UpdateInstanceLogParams{
		ID:              req.ID,
		Type:            req.Type,
		Title:           req.Title,
		Description:     req.Description,
		NullDescription: req.NullDescription,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, log)
}

type DeleteInstanceLogRequest struct {
	ID int64 `param:"id" validate:"required"`
}

func (h *EchoHandler) DeleteInstanceLog(c echo.Context) error {
	var req DeleteInstanceLogRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := h.service.DeleteInstanceLog(c.Request().Context(), req.ID); err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Instance log deleted successfully")
}
