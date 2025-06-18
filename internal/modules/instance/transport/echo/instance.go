package instanceecho

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	accountsvc "github.com/wagecloud/wagecloud-server/internal/modules/account/service"
	instancesvc "github.com/wagecloud/wagecloud-server/internal/modules/instance/service"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/shared/transport/http/response"
)

type EchoHandler struct {
	service instancesvc.Service
}

func NewEchoHandler(service instancesvc.Service) *EchoHandler {
	return &EchoHandler{service: service}
}

type GetInstanceRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) GetInstance(c echo.Context) error {
	var req GetInstanceRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	instance, err := h.service.GetInstance(c.Request().Context(), instancesvc.GetInstanceParams{
		ID: req.ID,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, instance)
}

func (h *EchoHandler) GetInstanceMonitor(c echo.Context) error {
	var req GetInstanceRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	monitor, err := h.service.GetInstanceMonitor(c.Request().Context(), req.ID)
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, monitor)
}

type ListInstancesRequest struct {
	Page          int32   `query:"page" validate:"min=1"`
	Limit         int32   `query:"limit" validate:"min=5,max=100"`
	NetworkID     *string `query:"network_id"`
	OsID          *string `query:"os_id"`
	ArchID        *string `query:"arch_id"`
	RegionID      *string `query:"region_id"`
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

func (h *EchoHandler) ListInstances(c echo.Context) error {
	var req ListInstancesRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusUnauthorized, err)
	}

	instances, err := h.service.ListInstances(c.Request().Context(), instancesvc.ListInstancesParams{
		PaginationParams: pagination.PaginationParams{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Account:       claims.ToAuthenticatedAccount(),
		OsID:          req.OsID,
		ArchID:        req.ArchID,
		RegionID:      req.RegionID,
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
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromPaginate(c.Response().Writer, instances)
}

type CreateInstanceRequest struct {
	Basic struct {
		Name     string `json:"name"`
		Hostname string `json:"hostname"`
		OsID     string `json:"os_id"`
		ArchID   string `json:"arch_id"`
		RegionID string `json:"region_id"`
	} `json:"basic"`
	Resources struct {
		Memory  int32 `json:"memory"`
		Cpu     int32 `json:"cpu"`
		Storage int32 `json:"storage"`
	} `json:"resources"`
	Security struct {
		Password          string   `json:"password"`
		SSHAuthorizedKeys []string `json:"ssh-authorized-keys"`
	} `json:"security"`
}

func (h *EchoHandler) CreateInstance(c echo.Context) error {
	var req CreateInstanceRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusUnauthorized, err)
	}

	paymentResult, err := h.service.PayCreateInstance(c.Request().Context(), instancesvc.PayCreateInstanceParams{
		CreateInstanceParams: instancesvc.CreateInstanceParams{
			Account:           claims.ToAuthenticatedAccount(),
			Name:              req.Basic.Name,
			SSHAuthorizedKeys: req.Security.SSHAuthorizedKeys,
			Password:          req.Security.Password,
			LocalHostname:     req.Basic.Hostname,
			OsID:              req.Basic.OsID,
			ArchID:            req.Basic.ArchID,
			RegionID:          req.Basic.RegionID,
			Memory:            req.Resources.Memory,
			Cpu:               req.Resources.Cpu,
			Storage:           req.Resources.Storage,
		},
		Method: paymentmodel.PaymentMethodVNPAY,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusCreated, struct {
		PaymentUrl string `json:"payment_url"`
		ID         int64  `json:"id,omitempty"`
	}{
		PaymentUrl: paymentResult.URL,
		ID:         paymentResult.Payment.ID, // ID will be set after payment is completed
	})
}

type UpdateInstanceRequest struct {
	ID        string  `param:"id" validate:"required,min=1,max=255"`
	NetworkID *string `json:"network_id"`
	OsID      *string `json:"os_id"`
	ArchID    *string `json:"arch_id"`
	Name      *string `json:"name"`
	Cpu       *int64  `json:"cpu"`
	Ram       *int64  `json:"ram"`
	Storage   *int64  `json:"storage"`
}

func (h *EchoHandler) UpdateInstance(c echo.Context) error {
	var req UpdateInstanceRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusUnauthorized, err)
	}

	instance, err := h.service.UpdateInstance(c.Request().Context(), instancesvc.UpdateInstanceParams{
		Account:   claims.ToAuthenticatedAccount(),
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
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromDTO(c.Response().Writer, http.StatusOK, instance)
}

type DeleteInstanceRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) DeleteInstance(c echo.Context) error {
	var req DeleteInstanceRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := h.service.DeleteInstance(c.Request().Context(), instancesvc.DeleteInstanceParams{
		ID: req.ID,
	}); err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Instance deleted successfully")
}

type StartInstanceRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) StartInstance(c echo.Context) error {
	var req StartInstanceRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	fmt.Println("Starting instance with ID:", req.ID)

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := h.service.StartInstance(c.Request().Context(), instancesvc.StartInstanceParams{
		ID: req.ID,
	}); err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Instance started successfully")
}

type StopInstanceRequest struct {
	ID string `param:"id" validate:"required,min=1,max=255"`
}

func (h *EchoHandler) StopInstance(c echo.Context) error {
	var req StopInstanceRequest
	if err := c.Bind(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	if err := c.Validate(&req); err != nil {
		return response.FromError(c.Response().Writer, http.StatusBadRequest, err)
	}

	claims, err := accountsvc.GetClaims(c.Request())
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusUnauthorized, err)
	}

	err = h.service.StopInstance(c.Request().Context(), instancesvc.StopInstanceParams{
		Account: claims.ToAuthenticatedAccount(),
		ID:      req.ID,
	})
	if err != nil {
		return response.FromError(c.Response().Writer, http.StatusInternalServerError, err)
	}

	return response.FromMessage(c.Response().Writer, http.StatusOK, "Instance stopped successfully")
}
