package vm

import (
	"context"

	"github.com/google/uuid"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/util/hash"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo    repository.Repository
	libvirt libvirt.ServiceInterface
}

type ServiceInterface interface {
	GetVM(ctx context.Context, params GetVMParams) (VM, error)
	ListVMs(ctx context.Context, params ListVMsParams) (model.PaginateResult[VM], error)
	CreateVM(ctx context.Context, params CreateVMParams) (VM, error)
	UpdateVM(ctx context.Context, params UpdateVMParams) (VM, error)
	DeleteVM(ctx context.Context, params DeleteVMParams) error
	StartVM(ctx context.Context, params StartVMParams) error
	StopVM(ctx context.Context, params StopVMParams) error
}

func NewService(repo repository.Repository, libvirt libvirt.ServiceInterface) *Service {
	return &Service{repo: repo, libvirt: libvirt}
}

type GetVMParams struct {
	Role      model.Role
	AccountID int64
	ID        string
}

func (s *Service) withStatus(vm model.VM) (VM, error) {
	isActive, err := s.libvirt.IsActive(vm.ID)
	if err != nil {
		return VM{}, err
	}

	var status VMStatus
	if isActive {
		status = VMStatusRunning
	} else {
		status = VMStatusStopped
	}

	return VM{VM: vm, Status: status}, nil
}

func (s *Service) GetVM(ctx context.Context, params GetVMParams) (VM, error) {

	repoParams := repository.GetVMParams{
		ID:        params.ID,
		AccountID: &params.AccountID,
	}

	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	vm, err := s.repo.GetVM(ctx, repoParams)
	if err != nil {
		return VM{}, err
	}

	return s.withStatus(vm)
}

type ListVMsParams struct {
	model.PaginationParams
	Role          model.Role
	AccountID     int64
	NetworkID     *string
	OsID          *string
	ArchID        *string
	Name          *string
	CpuFrom       *int64
	CpuTo         *int64
	RamFrom       *int64
	RamTo         *int64
	StorageFrom   *int64
	StorageTo     *int64
	CreatedAtFrom *int64
	CreatedAtTo   *int64
}

func (s *Service) ListVMs(ctx context.Context, params ListVMsParams) (res model.PaginateResult[VM], err error) {
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

	vmsWithStatus := make([]VM, len(vms))
	for i, vm := range vms {
		vmsWithStatus[i], err = s.withStatus(vm)
		if err != nil {
			return res, err
		}
	}

	return model.PaginateResult[VM]{
		Data:     vmsWithStatus,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

type CreateVMParams struct {
	AccountID int64
	// Userdata
	Name              string
	SSHAuthorizedKeys []string
	Password          string
	// Metadata
	LocalHostname string
	//Spec
	OsID    string
	ArchID  string
	Memory  int
	Cpu     int
	Storage int
}

func (s *Service) CreateVM(ctx context.Context, params CreateVMParams) (VM, error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return VM{}, err
	}
	defer txRepo.Rollback(ctx)

	// 1. Create records in database
	os, err := txRepo.GetOS(ctx, params.OsID)
	if err != nil {
		return VM{}, err
	}

	arch, err := txRepo.GetArch(ctx, params.ArchID)
	if err != nil {
		return VM{}, err
	}

	network, err := txRepo.CreateNetwork(ctx, model.Network{
		ID:        uuid.New().String(),
		PrivateIP: "",
	})
	if err != nil {
		return VM{}, err
	}

	vm, err := txRepo.CreateVM(ctx, model.VM{
		AccountID: params.AccountID,
		NetworkID: network.ID,
		OsID:      os.ID,
		ArchID:    arch.ID,
		Name:      params.Name,
		Cpu:       int32(params.Cpu),
		Ram:       int32(params.Memory),
		Storage:   int32(params.Storage),
	})
	if err != nil {
		return VM{}, err
	}

	// 2. Create cloudinit
	userdata := libvirt.NewDefaultUserdata()
	userdata.Users[0].Name = params.Name
	userdata.Users[0].SSHAuthorizedKeys = params.SSHAuthorizedKeys
	userdata.Users[0].Passwd, err = hash.Password(params.Password)
	if err != nil {
		return VM{}, err
	}

	metadata := libvirt.NewDefaultMetadata()
	metadata.LocalHostname = params.LocalHostname

	networkConfig := libvirt.NewDefaultNetworkConfig()

	domain := libvirt.FromVMToDomain(vm)

	if err = s.libvirt.CreateCloudinit(libvirt.CreateCloudinitParams{
		Filepath:      domain.CloudinitPath(),
		Userdata:      userdata,
		Metadata:      metadata,
		NetworkConfig: networkConfig,
	}); err != nil {
		return VM{}, err
	}

	// 3. Create domain
	if err := s.libvirt.CreateDomain(domain); err != nil {
		return VM{}, err
	}

	if err := txRepo.Commit(ctx); err != nil {
		return VM{}, err
	}

	return s.withStatus(vm)
}

type UpdateVMParams struct {
	Role      model.Role
	AccountID int64
	ID        string
	NetworkID *string
	OsID      *string
	ArchID    *string
	Name      *string
	Cpu       *int32
	Ram       *int32
	Storage   *int32
}

func (s *Service) UpdateVM(ctx context.Context, params UpdateVMParams) (VM, error) {
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
		return VM{}, err
	}

	return s.withStatus(updatedVM)
}

type DeleteVMParams struct {
	ID        string
	AccountID int64
	Role      model.Role
}

func (s *Service) DeleteVM(ctx context.Context, params DeleteVMParams) error {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer txRepo.Rollback(ctx)

	repoParams := repository.DeleteVMParams{
		ID: params.ID,
	}

	// Users can only see their own VMs
	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	if err := txRepo.DeleteVM(ctx, repoParams); err != nil {
		return err
	}

	if err := txRepo.Commit(ctx); err != nil {
		return err
	}

	// ! Delete domain does not support rollback operation so it should done last (after commit)
	if err := s.libvirt.DeleteDomain(params.ID); err != nil {
		return err
	}

	return nil
}

type StartVMParams struct {
	AccountID int64
	Role      model.Role
	ID        string
}

func (s *Service) StartVM(ctx context.Context, params StartVMParams) error {
	// Users can only start their own VMs
	if params.Role == model.RoleUser {
		_, err := s.repo.GetVM(ctx, repository.GetVMParams{
			ID:        params.ID,
			AccountID: &params.AccountID,
		})
		if err != nil {
			return err
		}
	}

	// TODO: put this in background, kinda slow 💀
	return s.libvirt.StartDomain(params.ID)
}

type StopVMParams struct {
	AccountID int64
	Role      model.Role
	ID        string
}

func (s *Service) StopVM(ctx context.Context, params StopVMParams) error {
	// Users can only stop their own VMs
	if params.Role == model.RoleUser {
		_, err := s.repo.GetVM(ctx, repository.GetVMParams{
			ID:        params.ID,
			AccountID: &params.AccountID,
		})
		if err != nil {
			return err
		}
	}

	return s.libvirt.StopDomain(params.ID)
}
