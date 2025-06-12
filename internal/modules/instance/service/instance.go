package instancesvc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wagecloud/wagecloud-server/internal/client/libvirt"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/hash"
)

type ServiceImpl struct {
	osSvc   ossvc.Service
	libvirt libvirt.Client
	storage *instancestorage.Storage
}

type Service interface {
	// Instance
	GetInstance(ctx context.Context, params GetInstanceParams) (instancemodel.Instance, error)
	ListInstances(ctx context.Context, params ListInstancesParams) (pagination.PaginateResult[instancemodel.Instance], error)
	CreateInstance(ctx context.Context, params CreateInstanceParams) (instancemodel.Instance, error)
	UpdateInstance(ctx context.Context, params UpdateInstanceParams) (instancemodel.Instance, error)
	DeleteInstance(ctx context.Context, params DeleteInstanceParams) error
	StartInstance(ctx context.Context, params StartInstanceParams) error
	StopInstance(ctx context.Context, params StopInstanceParams) error
	// RestartInstance(ctx context.Context, params RestartInstanceParams) error

	// Network
	GetNetwork(ctx context.Context, params GetNetworkParams) (instancemodel.Network, error)
	ListNetworks(ctx context.Context, params ListNetworksParams) (pagination.PaginateResult[instancemodel.Network], error)
	CreateNetwork(ctx context.Context, params CreateNetworkParams) (instancemodel.Network, error)
	UpdateNetwork(ctx context.Context, params UpdateNetworkParams) (instancemodel.Network, error)
	DeleteNetwork(ctx context.Context, params DeleteNetworkParams) error
}

func NewService(libvirt libvirt.Client, storage *instancestorage.Storage, osSvc ossvc.Service) Service {
	return &ServiceImpl{
		osSvc:   osSvc,
		libvirt: libvirt,
		storage: storage,
	}
}

type GetInstanceParams struct {
	Account accountmodel.AuthenticatedAccount
	ID      string
}

func (s *ServiceImpl) GetInstance(ctx context.Context, params GetInstanceParams) (instancemodel.Instance, error) {
	instance, err := s.storage.GetInstance(ctx, params.ID)
	if err != nil {
		return instancemodel.Instance{}, err
	}

	if err = s.canAccess(ctx, canAccessParams{
		Account:  params.Account,
		Instance: instance,
	}); err != nil {
		return instancemodel.Instance{}, err
	}

	return instance, nil
}

type ListInstancesParams struct {
	pagination.PaginationParams
	Account       accountmodel.AuthenticatedAccount
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

func (s *ServiceImpl) ListInstances(ctx context.Context, params ListInstancesParams) (res pagination.PaginateResult[instancemodel.Instance], err error) {
	storageParams := instancestorage.ListInstancesParams{
		PaginationParams: params.PaginationParams,
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

	// Authorization: users can only see their own instances
	if params.Account.Type == accountmodel.AccountTypeUser {
		storageParams.AccountID = &params.Account.AccountID
	}

	total, err := s.storage.CountInstances(ctx, storageParams)
	if err != nil {
		return res, err
	}

	instances, err := s.storage.ListInstances(ctx, storageParams)
	if err != nil {
		return res, err
	}

	return pagination.PaginateResult[instancemodel.Instance]{
		Data:     instances,
		Limit:    params.Limit,
		Page:     params.Page,
		Total:    total,
		NextPage: params.NextPage(total),
	}, nil
}

type CreateInstanceParams struct {
	Account accountmodel.AuthenticatedAccount
	// Userdata
	Name              string
	SSHAuthorizedKeys []string
	Password          string
	// Metadata
	LocalHostname string
	//Spec
	OsID    string
	ArchID  string
	Memory  int32
	Cpu     int32
	Storage int32
}

func (s *ServiceImpl) CreateInstance(ctx context.Context, params CreateInstanceParams) (instancemodel.Instance, error) {
	txStorage, err := s.storage.BeginTx(ctx)
	if err != nil {
		return instancemodel.Instance{}, err
	}
	defer txStorage.Rollback(ctx)

	// 1. Create records in database
	os, err := s.osSvc.GetOS(ctx, ossvc.GetOSParams{
		ID: params.OsID,
	})
	if err != nil {
		return instancemodel.Instance{}, err
	}

	arch, err := s.osSvc.GetArch(ctx, params.ArchID)
	if err != nil {
		return instancemodel.Instance{}, err
	}

	network, err := txStorage.CreateNetwork(ctx, instancemodel.Network{
		ID:        uuid.New().String(),
		PrivateIP: "",
	})
	if err != nil {
		return instancemodel.Instance{}, err
	}

	instance, err := txStorage.CreateInstance(ctx, instancemodel.Instance{
		ID:        uuid.New().String(),
		AccountID: params.Account.AccountID,
		NetworkID: network.ID,
		OSID:      os.ID,
		ArchID:    arch.ID,
		Name:      params.Name,
		CPU:       int32(params.Cpu),
		RAM:       int32(params.Memory),
		Storage:   int32(params.Storage),
	})
	if err != nil {
		return instancemodel.Instance{}, err
	}

	// 2. Create cloudinit
	userdata := libvirt.NewDefaultUserdata()
	userdata.Users[0].Name = params.Name
	userdata.Users[0].SSHAuthorizedKeys = params.SSHAuthorizedKeys
	userdata.Users[0].Passwd, err = hash.Password(params.Password)
	if err != nil {
		return instancemodel.Instance{}, err
	}

	metadata := libvirt.NewDefaultMetadata()
	metadata.LocalHostname = params.LocalHostname

	networkConfig := libvirt.NewDefaultNetworkConfig()

	// Convert from our model to libvirt Domain
	domain := libvirt.Domain{
		ID:     instance.ID,
		Name:   instance.Name,
		Memory: libvirt.Memory{Value: uint(instance.RAM), Unit: libvirt.UnitMB},
		Cpu:    libvirt.Cpu{Value: uint(instance.CPU)},
		OS: libvirt.OS{
			Name: os.ID,
			Type: "hvm",
			Arch: arch.ID,
		},
		Storage: uint(instance.Storage),
	}

	if err = s.libvirt.CreateCloudinit(ctx, libvirt.CreateCloudinitParams{
		Filepath:      domain.CloudinitPath(),
		Userdata:      userdata,
		Metadata:      metadata,
		NetworkConfig: networkConfig,
	}); err != nil {
		return instancemodel.Instance{}, err
	}

	// 3. Create domain
	if err := s.libvirt.CreateDomain(ctx, domain); err != nil {
		return instancemodel.Instance{}, err
	}

	if err := txStorage.Commit(ctx); err != nil {
		return instancemodel.Instance{}, err
	}

	return instance, nil
}

type UpdateInstanceParams struct {
	Account   accountmodel.AuthenticatedAccount
	ID        string
	NetworkID *string
	OsID      *string
	ArchID    *string
	Name      *string
	Cpu       *int64
	Ram       *int64
	Storage   *int64
}

func (s *ServiceImpl) UpdateInstance(ctx context.Context, params UpdateInstanceParams) (instancemodel.Instance, error) {
	storageParams := instancestorage.UpdateInstanceParams{
		ID:      params.ID,
		Name:    params.Name,
		CPU:     params.Cpu,
		RAM:     params.Ram,
		Storage: params.Storage,
	}

	// Users can only see their own instances
	if params.Account.Type == accountmodel.AccountTypeUser {
		storageParams.AccountID = &params.Account.AccountID
	}

	updatedInstance, err := s.storage.UpdateInstance(ctx, storageParams)
	if err != nil {
		return instancemodel.Instance{}, err
	}

	return updatedInstance, nil
}

type DeleteInstanceParams struct {
	Account accountmodel.AuthenticatedAccount
	ID      string
}

func (s *ServiceImpl) DeleteInstance(ctx context.Context, params DeleteInstanceParams) error {
	txStorage, err := s.storage.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer txStorage.Rollback(ctx)

	// TODO: missing checK: user only delete their own instances

	if err := txStorage.DeleteInstance(ctx, params.ID); err != nil {
		return err
	}

	if err := txStorage.Commit(ctx); err != nil {
		return err
	}

	// ! Delete domain does not support rollback operation so it should done last (after commit)
	// TODO: move this libvirt create/delete logic to storage to support atomic operation (?)
	if err := s.libvirt.DeleteDomain(ctx, params.ID); err != nil {
		return err
	}

	return nil
}

type StartInstanceParams struct {
	Account accountmodel.AuthenticatedAccount
	ID      string
}

func (s *ServiceImpl) StartInstance(ctx context.Context, params StartInstanceParams) error {
	// TODO: put this in background, kinda slow 💀
	return s.libvirt.StartDomain(ctx, params.ID)
}

type StopInstanceParams struct {
	Account accountmodel.AuthenticatedAccount
	ID      string
}

func (s *ServiceImpl) StopInstance(ctx context.Context, params StopInstanceParams) error {
	return s.libvirt.StopDomain(ctx, params.ID)
}

type canAccessParams struct {
	Account  accountmodel.AuthenticatedAccount
	Instance instancemodel.Instance
}

// TODO: future upgrade
func (s *ServiceImpl) canAccess(_ context.Context, params canAccessParams) error {
	// Users can only access their own instances
	if params.Account.Type == accountmodel.AccountTypeUser && params.Account.AccountID != params.Instance.AccountID {
		return errors.New("access denied: user can only access their own instances")
	}

	if params.Account.Type == accountmodel.AccountTypeAdmin {
		return nil
	}

	return errors.New("access denied: unsupported role or instance access")
}
