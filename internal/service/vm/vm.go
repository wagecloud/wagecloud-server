package vm

import (
	"context"

	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo    repository.Repository
	libvirt libvirt.ServiceInterface
}

type ServiceInterface interface {
	GetVM(ctx context.Context, params GetVMParams) (model.VM, error)
	ListVMs(ctx context.Context, params ListVMsParams) (model.PaginateResult[model.VM], error)
	CreateVM(ctx context.Context, params CreateVMParams) (model.VM, error)
	UpdateVM(ctx context.Context, params UpdateVMParams) (model.VM, error)
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

func (s *Service) GetVM(ctx context.Context, params GetVMParams) (model.VM, error) {

	repoParams := repository.GetVMParams{
		ID:        params.ID,
		AccountID: &params.AccountID,
	}

	if params.Role == model.RoleUser {
		repoParams.AccountID = &params.AccountID
	}

	return s.repo.GetVM(ctx, repoParams)
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

func (s *Service) CreateVM(ctx context.Context, params CreateVMParams) (model.VM, error) {
	txRepo, err := s.repo.Begin(ctx)
	if err != nil {
		return model.VM{}, err
	}
	defer txRepo.Rollback(ctx)

	// 1. Create records in database
	os, err := txRepo.GetOS(ctx, params.OsID)
	if err != nil {
		return model.VM{}, err
	}

	arch, err := txRepo.GetArch(ctx, params.ArchID)
	if err != nil {
		return model.VM{}, err
	}

	network, err := txRepo.CreateNetwork(ctx, model.Network{
		PrivateIP: "",
	})
	if err != nil {
		return model.VM{}, err
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
		return model.VM{}, err
	}

	// 2. Create cloudinit
	userdata := libvirt.NewDefaultUserdata()
	userdata.Users[0].Name = params.Name
	userdata.Users[0].SSHAuthorizedKeys = params.SSHAuthorizedKeys
	userdata.Users[0].Passwd = params.Password

	metadata := libvirt.NewDefaultMetadata()
	metadata.LocalHostname = params.LocalHostname

	networkConfig := libvirt.NewDefaultNetworkConfig()

	domain := libvirt.ToDomain(vm)

	if err = s.libvirt.CreateCloudinit(libvirt.CreateCloudinitParams{
		Filepath:      domain.CloudinitPath(),
		Userdata:      userdata,
		Metadata:      metadata,
		NetworkConfig: networkConfig,
	}); err != nil {
		return model.VM{}, err
	}

	// 3. Create domain
	if err := s.libvirt.CreateDomain(domain); err != nil {
		return model.VM{}, err
	}

	if err := txRepo.Commit(ctx); err != nil {
		return model.VM{}, err
	}

	return vm, nil
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
	ID        string
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

type StartVMParams struct {
	AccountID int64
	Role      model.Role
	ID        string
}

func (s *Service) StartVM(ctx context.Context, params StartVMParams) error {
	return nil
}

type StopVMParams struct {
	AccountID int64
	Role      model.Role
	ID        string
}

func (s *Service) StopVM(ctx context.Context, params StopVMParams) error {
	return nil
}
