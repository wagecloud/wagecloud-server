package instancesvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/wagecloud/wagecloud-server/internal/client/libvirt"
	"github.com/wagecloud/wagecloud-server/internal/client/nats"
	"github.com/wagecloud/wagecloud-server/internal/client/redis"
	"github.com/wagecloud/wagecloud-server/internal/logger"
	accountmodel "github.com/wagecloud/wagecloud-server/internal/modules/account/model"
	instancemodel "github.com/wagecloud/wagecloud-server/internal/modules/instance/model"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	ossvc "github.com/wagecloud/wagecloud-server/internal/modules/os/service"
	paymentmodel "github.com/wagecloud/wagecloud-server/internal/modules/payment/model"
	paymentsvc "github.com/wagecloud/wagecloud-server/internal/modules/payment/service"
	commonmodel "github.com/wagecloud/wagecloud-server/internal/shared/model"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
	"github.com/wagecloud/wagecloud-server/internal/utils/hash"
)

type ServiceImpl struct {
	storage    *instancestorage.Storage
	redis      redis.Client
	nats       nats.Client
	libvirt    libvirt.Client
	osSvc      ossvc.Service
	paymentSvc paymentsvc.Service
	cron       *cron.Cron
}

type Service interface {
	// Instance
	GetInstance(ctx context.Context, params GetInstanceParams) (instancemodel.Instance, error)
	ListInstances(ctx context.Context, params ListInstancesParams) (pagination.PaginateResult[instancemodel.Instance], error)
	CreateInstance(ctx context.Context, params CreateInstanceParams) (instancemodel.Instance, error)
	PayCreateInstance(ctx context.Context, params PayCreateInstanceParams) (PayCreateInstanceResult, error)
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
	MapPortNginx(ctx context.Context, params MapPortNginxParams) error
	UnmapPortNginx(ctx context.Context, params UnmapPortNginxParams) error

	// Domain
	GetDomain(ctx context.Context, id int64) (instancemodel.Domain, error)
	ListDomains(ctx context.Context, params ListDomainsParams) (pagination.PaginateResult[instancemodel.Domain], error)
	CreateDomain(ctx context.Context, params CreateDomainParams) (instancemodel.Domain, error)
	UpdateDomain(ctx context.Context, params UpdateDomainParams) (instancemodel.Domain, error)
	DeleteDomain(ctx context.Context, id int64) error

	// Instance Log
	GetInstanceLog(ctx context.Context, id int64) (instancemodel.InstanceLog, error)
	ListInstanceLogs(ctx context.Context, params ListInstanceLogsParams) (pagination.PaginateResult[instancemodel.InstanceLog], error)
	CreateInstanceLog(ctx context.Context, params CreateInstanceLogParams) (instancemodel.InstanceLog, error)
	UpdateInstanceLog(ctx context.Context, params UpdateInstanceLogParams) (instancemodel.InstanceLog, error)
	DeleteInstanceLog(ctx context.Context, id int64) error

	// Region
	GetRegion(ctx context.Context, id string) (instancemodel.Region, error)
	ListRegions(ctx context.Context, params ListRegionsParams) (pagination.PaginateResult[instancemodel.Region], error)
	CreateRegion(ctx context.Context, params CreateRegionParams) (instancemodel.Region, error)
	UpdateRegion(ctx context.Context, params UpdateRegionParams) (instancemodel.Region, error)
	DeleteRegion(ctx context.Context, id string) error
}

func NewService(libvirt libvirt.Client, nats nats.Client, redis redis.Client, storage *instancestorage.Storage, osSvc ossvc.Service, paymentSvc paymentsvc.Service) Service {
	s := &ServiceImpl{
		nats:       nats,
		redis:      redis,
		osSvc:      osSvc,
		libvirt:    libvirt,
		storage:    storage,
		paymentSvc: paymentSvc,
		cron:       cron.New(cron.WithSeconds()),
	}
	s.init()

	// TODO: shitty ass code, refactor later
	s.cron.AddFunc("@every 30s", func() {
		s.UpdateNetworkIPs(context.Background())
	})
	s.cron.Start()

	// s.UpdateNetworkIPs(context.Background())

	return s
}

func (s *ServiceImpl) init() {
	// TODO: refactor shitass code
	s.nats.Subscribe("payment.processed", func(data []byte) {
		ctx := context.Background()
		var paymentNAT paymentmodel.PaymentProcesseDataNATS
		if err := json.Unmarshal(data, &paymentNAT); err != nil {
			logger.Log.Error("failed to unmarshal payment processed data from NATS: " + err.Error())
			return
		}

		logger.Log.Info(fmt.Sprintf("received payment processed event: %+v", paymentNAT))

		redisKey := "pay_create_instance:" + strconv.FormatInt(paymentNAT.PaymentID, 10)

		byteData, err := s.redis.Get(ctx, redisKey)
		if err != nil {
			logger.Log.Error("failed to get payment data from Redis: " + err.Error())
		}

		var params CreateInstanceParams
		if err := json.Unmarshal(byteData, &params); err != nil {
			logger.Log.Error("failed to unmarshal payment data: " + err.Error())
			return
		}

		fmt.Println("Creating instance after payment processed with params: ", params)

		instance, err := s.CreateInstance(ctx, params)
		if err != nil {
			logger.Log.Error("failed to create instance after payment: " + err.Error())
			return
		}

		if err := s.redis.Delete(ctx, redisKey); err != nil {
			logger.Log.Error("failed to delete payment data from Redis: " + err.Error())
		} else {
			logger.Log.Info("successfully deleted payment data from Redis after creating instance")
		}

		fmt.Printf("successfully created instance %s after payment processed: %+v \n", instance.ID, params)
	})
}

func (s *ServiceImpl) UpdateNetworkIPs(ctx context.Context) {
	domains, err := s.libvirt.ListDomains(ctx, libvirt.ListDomainsParams{
		Flags: 16, // libvirt.CONNECT_LIST_DOMAINS_RUNNING,
	})
	if err != nil {
		logger.Log.Error("failed to list domains: " + err.Error())
		return
	}

	for _, domain := range domains {
		ip, err := s.libvirt.GetPrivateIP(ctx, domain.ID)
		if err != nil {
			fmt.Println("failed to get private IP for domain", domain.ID, ":", err)
			continue
		}

		fmt.Printf("Domain %s has private IP: %s\n", domain.ID, ip)

		// Update the network IP in the database
		if _, err := s.storage.UpdateNetwork(ctx, instancestorage.UpdateNetworkParams{
			InstanceID: &domain.ID,
			PrivateIP:  &ip,
		}); err != nil {
			fmt.Println("failed to update network IP for domain", domain.ID, ":", err)
			continue
		}

		fmt.Printf("Successfully updated network IP for domain %s to %s\n", domain.ID, ip)
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

	instanceID := uuid.New().String()

	instance, err := txStorage.CreateInstance(ctx, instancemodel.Instance{
		ID:        instanceID,
		AccountID: params.Account.AccountID,
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

	_, err = txStorage.CreateNetwork(ctx, instancemodel.Network{
		InstanceID: instance.ID,
		PrivateIP:  "",
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

type PayCreateInstanceParams struct {
	CreateInstanceParams
	Method paymentmodel.PaymentMethod
}

type PayCreateInstanceResult struct {
	Payment paymentmodel.Payment
	Items   []paymentmodel.PaymentItem
	URL     string
}

// PayAndCreateInstance creates a new instance and returns the payment URL for the user to pay.
// Waits for the payment to be successful before creating the instance.
func (s *ServiceImpl) PayCreateInstance(ctx context.Context, params PayCreateInstanceParams) (PayCreateInstanceResult, error) {
	// TODO: remove hard-coded example price:
	// Storage: 100.000 VND/GB
	// Memory: 150.000 VND/GB
	// CPU: 200.000 VND/CPU

	totalPrice := commonmodel.NewConcurrency(float64(params.Storage)*100_000) +
		commonmodel.NewConcurrency(float64(params.Memory/1024)*150_000) +
		commonmodel.NewConcurrency(float64(params.Cpu)*200_000)

	paymentResult, err := s.paymentSvc.CreatePayment(ctx, paymentsvc.CreatePaymentParams{
		Account: params.Account,
		Method:  params.Method,
		Items: []paymentsvc.CreatePaymentParamsItem{{
			Name:  params.Name,
			Price: totalPrice,
		}},
	})
	if err != nil {
		return PayCreateInstanceResult{}, err
	}

	byteData, err := json.Marshal(params.CreateInstanceParams)
	if err != nil {
		return PayCreateInstanceResult{}, fmt.Errorf("failed to marshal payment data: %w", err)
	}

	if err = s.redis.Set(ctx, "pay_create_instance:"+strconv.FormatInt(paymentResult.Payment.ID, 10), byteData, 5*time.Minute); err != nil {
		return PayCreateInstanceResult{}, fmt.Errorf("failed to set payment data in Redis: %w", err)
	}

	return PayCreateInstanceResult{
		Payment: paymentResult.Payment,
		Items:   paymentResult.Items,
		URL:     paymentResult.URL,
	}, nil
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
