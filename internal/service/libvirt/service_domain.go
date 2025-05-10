package libvirt

import (
	"errors"
	"fmt"
	"io"
	"os"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/wagecloud/wagecloud-server/internal/logger"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/util/file"
	"github.com/wagecloud/wagecloud-server/internal/util/transaction"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo    *repository.RepositoryImpl
	connect *libvirt.Connect
	tx      *transaction.Transaction
}

type ServiceInterface interface {
	WithTx(tx *transaction.Transaction) ServiceInterface

	// CLOUDINIT
	CreateCloudinit(params CreateCloudinitParams) error
	CreateCloudinitByReader(params CreateCloudinitByReaderParams) error
	WriteCloudinit(userdata io.Reader, metadata io.Reader, networkConfig io.Reader, cloudinitFile io.Writer) error

	// DOMAIN
	GetDomain(domainID string) (Domain, error)
	IsActive(domainID string) (bool, error)
	ListDomains(params ListDomainsParams) ([]Domain, error)
	CreateDomain(domain Domain) error
	UpdateDomain(domainID string, params UpdateDomainParams) error
	DeleteDomain(domainID string) error
	StartDomain(domainID string) error
	StopDomain(domainID string) error

	// QEMU
	CreateImage(params CreateImageParams) error
}

const (
	QemuConnect = "qemu:///system"
)

var (
	ErrDomainNotFound = errors.New("domain not found")
)

func NewService(repo *repository.RepositoryImpl) *Service {
	return &Service{repo: repo, tx: transaction.NewTransaction(true)}
}

func (s *Service) getConnect() (*libvirt.Connect, error) {
	if s.connect == nil {
		conn, err := libvirt.NewConnect(QemuConnect)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to libvirt: %v", err)
		}
		s.connect = conn
	}
	return s.connect, nil
}

func (s *Service) getDomain(domainID string) (*libvirt.Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return nil, err
	}

	domain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return nil, ErrDomainNotFound
	}

	return domain, nil
}

func (s *Service) WithTx(tx *transaction.Transaction) ServiceInterface {
	newSvc := NewService(s.repo)
	newSvc.tx = tx
	return newSvc
}

func (s *Service) GetDomain(domainID string) (Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return Domain{}, err
	}

	domain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return Domain{}, ErrDomainNotFound
	}

	return FromLibvirtToDomain(*domain)
}

func (s *Service) IsActive(domainID string) (bool, error) {
	domain, err := s.getDomain(domainID)
	if err != nil {
		return false, err
	}
	return domain.IsActive()
}

type ListDomainsParams struct {
	Flags libvirt.ConnectListAllDomainsFlags
}

func (s *Service) ListDomains(params ListDomainsParams) ([]Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return nil, err
	}

	domains, err := conn.ListAllDomains(params.Flags)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %v", err)
	}

	domainsModel := make([]Domain, len(domains))

	for i, domain := range domains {
		model, err := FromLibvirtToDomain(domain)
		if err != nil {
			return nil, fmt.Errorf("failed to convert domain to model: %v", err)
		}

		domainsModel[i] = model
	}

	return domainsModel, nil
}

// CreateDomain creates a new domain in libvirt
//
// Supports rollback operation, safe to use in anywhere, anytime
func (s *Service) CreateDomain(domain Domain) error {
	conn, err := s.getConnect()
	if err != nil {
		return err
	}

	tx := transaction.NewTransaction(false)
	svcTx := s.WithTx(tx)
	defer tx.Rollback()

	// Create new qcow2 image from base image
	if err = svcTx.CreateImage(CreateImageParams{
		BaseImagePath:  domain.BaseImagePath(),
		CloneImagePath: domain.VMImagePath(),
		Size:           domain.Storage,
	}); err != nil {
		return fmt.Errorf("failed to clone image: %v", err)
	}

	domainXML, err := getXMLConfig(domain)
	if err != nil {
		return fmt.Errorf("failed to generate domain XML: %v", err)
	}

	xmlData, err := domainXML.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal domain XML: %v", err)
	}

	_, err = conn.DomainDefineXML(xmlData)
	if err != nil {
		return fmt.Errorf("failed to define domain: %v", err)
	}

	tx.Commit()

	return nil
}

type UpdateDomainParams struct {
	Name    *string
	Cpu     *uint
	Ram     *uint
	Storage *uint
}

func (s *Service) UpdateDomain(domainID string, params UpdateDomainParams) error {
	conn, err := s.getConnect()
	if err != nil {
		return err
	}

	libDomain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	domain, err := FromLibvirtToDomain(*libDomain)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %v", err)
	}

	// Start updating domain
	if params.Name != nil {
		domain.Name = *params.Name
	}
	if params.Cpu != nil {
		domain.Cpu.Value = *params.Cpu
	}
	if params.Ram != nil {
		domain.Memory.Value = *params.Ram
	}
	if params.Storage != nil {
		domain.Storage = *params.Storage
	}

	// Generate new XML config
	domainXML, err := getXMLConfig(domain)
	if err != nil {
		return fmt.Errorf("failed to generate domain XML: %v", err)
	}

	xmlData, err := domainXML.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal domain XML: %v", err)
	}

	_, err = conn.DomainDefineXML(xmlData)
	if err != nil {
		return fmt.Errorf("failed to define domain: %v", err)
	}

	// TODO: real update domain, not redefine and start it again
	// If the domain is running, we need to start it again
	if val, err := libDomain.IsActive(); err != nil && val {
		if err := libDomain.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy domain: %v", err)
		}
	}

	return nil
}

// DeleteDomain removes the domain from libvirt
//
// Does not support rollback operation so it should done last
func (s *Service) DeleteDomain(domainID string) error {
	libDomain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	domain, err := FromLibvirtToDomain(*libDomain)
	if err != nil {
		return fmt.Errorf("failed to convert domain to model: %v", err)
	}

	isActive, err := libDomain.IsActive()
	if err != nil {
		return fmt.Errorf("failed to check if domain is active: %v", err)
	}

	if isActive {
		if err := libDomain.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy domain: %v", err)
		}
	}

	if err = libDomain.Undefine(); err != nil {
		return fmt.Errorf("failed to undefine domain: %v", err)
	}

	// remove vm and cloudinit (always after domain is stopped and should not return error)
	//! These removal operations cannot be rolled back, so it should be done last
	if err := os.Remove(domain.VMImagePath()); err != nil {
		logger.Log.Error("failed to remove vm image", zap.String("path", domain.VMImagePath()), zap.Error(err))
	}
	if err := os.Remove(domain.CloudinitPath()); err != nil {
		logger.Log.Error("failed to remove cloudinit", zap.String("path", domain.CloudinitPath()), zap.Error(err))
	}

	return nil
}

func (s *Service) StartDomain(domainID string) error {
	domain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	return domain.Create()
}

func (s *Service) StopDomain(domainID string) error {
	domain, err := s.getDomain(domainID)
	if err != nil {
		return err
	}

	return domain.Shutdown()
}

func getXMLConfig(domain Domain) (*libvirtxml.Domain, error) {
	vmImagePath := domain.VMImagePath()
	cloudinitPath := domain.CloudinitPath()

	if !file.Exists(vmImagePath) || !file.Exists(cloudinitPath) {
		return nil, fmt.Errorf("image or cloudinit file not found")
	}

	domainXML := &libvirtxml.Domain{
		Type: "kvm",
		Name: domain.ID,
		UUID: domain.ID,
		Memory: &libvirtxml.DomainMemory{
			Value: domain.Memory.Value,
			Unit:  string(domain.Memory.Unit),
		},
		CurrentMemory: &libvirtxml.DomainCurrentMemory{
			Value: domain.Memory.Value,
			Unit:  string(domain.Memory.Unit),
		},
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			Value:     domain.Cpu.Value,
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    domain.OS.Arch, // x86_64
				Machine: "pc-q35-6.2",
				Type:    domain.OS.Type, // hvm
			},
		},
		CPU: &libvirtxml.DomainCPU{
			Mode:       "host-passthrough",
			Check:      "none",
			Migratable: "on",
		},
		Clock: &libvirtxml.DomainClock{
			Offset: "utc",
		},
		OnPoweroff: "destroy",
		OnReboot:   "destroy",
		OnCrash:    "destroy",
		Devices: &libvirtxml.DomainDeviceList{
			Disks: []libvirtxml.DomainDisk{
				{
					Device: "disk",
					Driver: &libvirtxml.DomainDiskDriver{
						Name:    "qemu",
						Type:    "qcow2",
						Discard: "unmap",
					},
					Source: &libvirtxml.DomainDiskSource{
						File: &libvirtxml.DomainDiskSourceFile{
							File: vmImagePath,
						},
					},
					Target: &libvirtxml.DomainDiskTarget{
						Dev: "vda",
						Bus: "virtio",
					},
				},
				{
					Device: "cdrom",
					Driver: &libvirtxml.DomainDiskDriver{
						Name: "qemu",
						Type: "raw",
					},
					Source: &libvirtxml.DomainDiskSource{
						File: &libvirtxml.DomainDiskSourceFile{
							File: cloudinitPath,
						},
					},
					Target: &libvirtxml.DomainDiskTarget{
						Dev: "sdb",
						Bus: "sata",
					},
					ReadOnly: &libvirtxml.DomainDiskReadOnly{},
				},
			},
			Interfaces: []libvirtxml.DomainInterface{
				{
					MAC: &libvirtxml.DomainInterfaceMAC{
						Address: "52:54:00:b7:a5:c2",
					},
					Source: &libvirtxml.DomainInterfaceSource{
						Bridge: &libvirtxml.DomainInterfaceSourceBridge{
							Bridge: "virbr0",
						},
					},
					Model: &libvirtxml.DomainInterfaceModel{
						Type: "virtio",
					},
				},
			},
			Graphics: []libvirtxml.DomainGraphic{
				{
					VNC: &libvirtxml.DomainGraphicVNC{
						Port:   -1,
						Listen: "0.0.0.0",
					},
				},
			},
			Consoles: []libvirtxml.DomainConsole{
				{
					TTY: "pty",
					Target: &libvirtxml.DomainConsoleTarget{
						Type: "serial",
					},
				},
			},
		},
	}

	return domainXML, nil
}
