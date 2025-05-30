package instanceecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

type EchoHandler struct {
	instanceSvc instancesvc.Service
}

func NewEchoHandler(instanceSvc instancesvc.Service) *EchoHandler {
	return &EchoHandler{
		instanceSvc: instanceSvc,
	}
}

func (h *EchoHandler) GetInstance(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	Instance, err := h.instanceSvc.GetInstance(c.Request().Context(), instancesvc.GetInstanceParams{
		Role:      claims.Role,
		AccountID: claims.AccountID,
		ID:        req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Instance)
}

func (h *EchoHandler) ListInstances(c echo.Context) error {
	var req struct {
		Page          int32   `query:"page" validate:"required,min=1"`
		Limit         int32   `query:"limit" validate:"required,min=5,max=100"`
		NetworkID     *string `query:"network_id"`
		OsID          *string `query:"os_id"`
		ArchID        *string `query:"arch_id"`
		Name          *string `query:"name"`
		CpuFrom       *int64  `query:"cpu_from"`
		CpuTo         *int64  `query:"cpu_to"`
		RamFrom       *int64  `query:"ram_from"`
		RamTo         *int64  `query:"ram_to"`
		StorageFrom   *int64  `query:"storage_from"`
		StorageTo     *int64  `query:"storage_to"`
		CreatedAtFrom *int64  `query:"created_at_from"`
		CreatedAtTo   *int64  `query:"created_at_to"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	Instances, err := h.instanceSvc.ListInstances(c.Request().Context(), instancesvc.ListInstancesParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		AccountID:     claims.AccountID,
		Role:          claims.Role,
		NetworkID:     req.NetworkID,
		OsID:          req.OsID,
		ArchID:        req.ArchID,
		Name:          req.Name,
		CpuFrom:       req.CpuFrom,
		CpuTo:         req.CpuTo,
		RamFrom:       req.RamFrom,
		RamTo:         req.RamTo,
		StorageFrom:   req.StorageFrom,
		StorageTo:     req.StorageTo,
		CreatedAtFrom: req.CreatedAtFrom,
		CreatedAtTo:   req.CreatedAtTo,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Instances)
}

func (h *EchoHandler) CreateInstance(c echo.Context) error {
	var req struct {
		Userdata struct {
			Name              string   `json:"name" validate:"required,min=1,max=255"`
			SSHAuthorizedKeys []string `json:"ssh-authorized-keys" validate:"omitempty,dive,min=20,max=5000"`
			Password          string   `json:"password" validate:"required,min=8,max=72"`
		} `json:"userdata" validate:"required"`
		Metadata struct {
			LocalHostname string `json:"local-hostname" validate:"required,hostname"`
		} `json:"metadata" validate:"required"`
		Spec struct {
			OsID    string `json:"os_id" validate:"required,min=1,max=255"`
			ArchID  string `json:"arch_id" validate:"required,min=1,max=255"`
			Memory  int    `json:"memory" validate:"required,min=512,max=262144"`
			Cpu     int    `json:"cpu" validate:"required,min=1,max=64"`
			Storage int    `json:"storage" validate:"required,min=10,max=2048"`
		} `json:"spec" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	Instance, err := h.instanceSvc.CreateInstance(c.Request().Context(), instancesvc.CreateInstanceParams{
		AccountID:         claims.AccountID,
		Name:              req.Userdata.Name,
		SSHAuthorizedKeys: req.Userdata.SSHAuthorizedKeys,
		Password:          req.Userdata.Password,
		LocalHostname:     req.Metadata.LocalHostname,
		OsID:              req.Spec.OsID,
		ArchID:            req.Spec.ArchID,
		Memory:            req.Spec.Memory,
		Cpu:               req.Spec.Cpu,
		Storage:           req.Spec.Storage,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Instance)
}

func (h *EchoHandler) UpdateInstance(c echo.Context) error {
	var req struct {
		ID        string  `param:"id" validate:"required,min=1,max=255"`
		NetworkID *string `json:"network_id"`
		OsID      *string `json:"os_id"`
		ArchID    *string `json:"arch_id"`
		Name      *string `json:"name"`
		Cpu       *int32  `json:"cpu"`
		Ram       *int32  `json:"ram"`
		Storage   *int32  `json:"storage"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	Instance, err := h.instanceSvc.UpdateInstance(c.Request().Context(), instancesvc.UpdateInstanceParams{
		Role:      claims.Role,
		AccountID: claims.AccountID,
		ID:        req.ID,
		NetworkID: req.NetworkID,
		OsID:      req.OsID,
		ArchID:    req.ArchID,
		Name:      req.Name,
		Cpu:       req.Cpu,
		Ram:       req.Ram,
		Storage:   req.Storage,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Instance)
}

func (h *EchoHandler) DeleteInstance(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.instanceSvc.DeleteInstance(c.Request().Context(), instancesvc.DeleteInstanceParams{
		Role:      claims.Role,
		AccountID: claims.AccountID,
		ID:        req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *EchoHandler) StartInstance(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.instanceSvc.StartInstance(c.Request().Context(), instancesvc.StartInstanceParams{
		AccountID: claims.AccountID,
		Role:      claims.Role,
		ID:        req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *EchoHandler) StopInstance(c echo.Context) error {
	var req struct {
		ID string `param:"id" validate:"required,min=1,max=255"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	err = h.instanceSvc.StopInstance(c.Request().Context(), instancesvc.StopInstanceParams{
		AccountID: claims.AccountID,
		Role:      claims.Role,
		ID:        req.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
