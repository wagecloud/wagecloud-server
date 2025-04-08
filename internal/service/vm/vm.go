package vm

import (
	"context"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo repository.Repository
}

type ServiceInterface interface {
	GetVM(ctx context.Context, params GetVMParams) (model.VM, error)
	ListVMs(ctx context.Context, params ListVMsParams) (model.PaginateResult[model.VM], error)
	CreateVM(ctx context.Context, vm model.VM) (model.VM, error)
	UpdateVM(ctx context.Context, params UpdateVMParams) (model.VM, error)
	DeleteVM(ctx context.Context, params DeleteVMParams) error
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

type GetVMParams struct {
	AccountID int64
	ID        int64
}

func (s *Service) GetVM(ctx context.Context, params GetVMParams) (model.VM, error) {
	return s.repo.GetVM(ctx, repository.GetVMParams{
		ID:        params.ID,
		AccountID: &params.AccountID,
	})
}

type ListVMsParams struct {
	model.PaginationParams
	AccountID     int64
	Role          model.Role
	NetworkID     *string
	OsID          *string
	ArchID        *string
	Name          *string
	CpuFrom       *int64
	CpuTo         *int64
	RamFrom       *int32
	RamTo         *int32
	StorageFrom   *int32
	StorageTo     *int32
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *Service) ListVMs(ctx context.Context, params ListVMsParams) (res model.PaginateResult[model.VM], err error) {
	repoParams := repository.ListVMsParams{
		PaginationParams: params.PaginationParams,
		AccountID:        &params.AccountID,
		NetworkID:        params.NetworkID,
		OsID:             params.OsID,
		ArchID:           params.ArchID,
		Name:             params.Name,
		CpuFrom:          params.CpuFrom,
		CpuTo:            params.CpuTo,
		RamFrom:          params.RamFrom,
		RamTo:            params.RamTo,
		StorageFrom:      params.StorageFrom,
		StorageTo:        params.StorageTo,
		CreatedAtFrom:    params.CreatedAtFrom,
		CreatedAtTo:      params.CreatedAtTo,
	}

	// Users can only see their own VMs
	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	total, err := s.repo.CountVMs(ctx, repoParams)
	if err != nil {
		return res, err
	}

	vms, err := s.repo.ListVMs(ctx, repoParams)
	if err != nil {
		return res, err
	}

	return model.PaginateResult[model.VM]{
		Data:     vms,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

func (s *Service) CreateVM(ctx context.Context, vm model.VM) (model.VM, error) {
	return s.repo.CreateVM(ctx, vm)
}

type UpdateVMParams struct {
	ID        int64
	AccountID int64
	Role      model.Role
	NetworkID *string
	OsID      *string
	ArchID    *string
	Name      *string
	Cpu       *int32
	Ram       *int32
	Storage   *int32
}

func (s *Service) UpdateVM(ctx context.Context, params UpdateVMParams) (model.VM, error) {
	repoParams := repository.UpdateVMParams{
		ID:        params.ID,
		NetworkID: params.NetworkID,
		OsID:      params.OsID,
		ArchID:    params.ArchID,
		Name:      params.Name,
		Cpu:       params.Cpu,
		Ram:       params.Ram,
		Storage:   params.Storage,
	}

	// Users can only see their own VMs
	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	updatedVM, err := s.repo.UpdateVM(ctx, repoParams)
	if err != nil {
		return model.VM{}, err
	}

	return updatedVM, nil
}

type DeleteVMParams struct {
	ID        int64
	AccountID int64
	Role      model.Role
}

func (s *Service) DeleteVM(ctx context.Context, params DeleteVMParams) error {
	repoParams := repository.DeleteVMParams{
		ID: params.ID,
	}

	// Users can only see their own VMs
	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	return s.repo.DeleteVM(ctx, repoParams)
}
