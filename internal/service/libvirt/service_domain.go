package libvirt

import (
	"errors"
	"fmt"
	"io"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/util/file"
	"github.com/wagecloud/wagecloud-server/internal/util/qemu"
	"libvirt.org/go/libvirt"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo    *repository.RepositoryImpl
	connect *libvirt.Connect
}

type ServiceInterface interface {
	// CLOUDINIT
	CreateCloudinit(params CreateCloudinitParams) error
	CreateCloudinitByReader(params CreateCloudinitByReaderParams) error
	WriteCloudinit(userdata io.Reader, metadata io.Reader, networkConfig io.Reader, cloudinitFile io.Writer) error

	// DOMAIN
	GetDomain(domainID string) (Domain, error)
	GetListDomains() ([]Domain, error)
	CreateDomain(domain Domain) error
	UpdateDomain(domainID string, domain Domain) error
	StartDomain(domainID string) error
	StopDomain(domainID string) error
}

const (
	QemuConnect = "qemu:///system"
)

var (
	ErrDomainNotFound = errors.New("domain not found")
)

func NewService(repo *repository.RepositoryImpl) *Service {
	return &Service{repo: repo}
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

func (s *Service) GetDomain(domainID string) (Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return Domain{}, err
	}

	domain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return Domain{}, ErrDomainNotFound
	}

	return toEntity(domain)
}

func (s *Service) GetListDomains() ([]Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return nil, err
	}

	domains, err := conn.ListAllDomains(0)

	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %v", err)
	}

	domainsModel := make([]Domain, len(domains))

	for i, domain := range domains {
		model, err := toEntity(&domain)
		if err != nil {
			return nil, fmt.Errorf("failed to convert domain to model: %v", err)
		}
		domainsModel[i] = model

	}

	return domainsModel, nil
}

func (s *Service) CreateDomain(domain Domain) error {
	conn, err := s.getConnect()
	if err != nil {
		return err
	}

	// Create new qcow2 image from base image
	if err = qemu.CreateImageWithPath(
		domain.BaseImagePath(),
		domain.VMImagePath(),
		domain.Storage,
	); err != nil {
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

	return nil

}

func (s *Service) UpdateDomain(domainID string, domain Domain) error {
	conn, err := s.getConnect()
	if err != nil {
		return err
	}

	domainObj, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return ErrDomainNotFound
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

	// If the domain is running, we need to start it again
	if val, err := domainObj.IsActive(); err != nil && val {
		if err := domainObj.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy domain: %v", err)
		}
	}

	return nil
}

func (s *Service) GetListActiveDomains() ([]Domain, error) {
	conn, err := s.getConnect()
	if err != nil {
		return nil, err
	}

	domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %v", err)
	}

	domainsModel := make([]Domain, len(domains))

	for i, domain := range domains {
		model, err := toEntity(&domain)
		if err != nil {
			return nil, fmt.Errorf("failed to convert domain to model: %v", err)
		}
		domainsModel[i] = model
	}

	return domainsModel, nil
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
		Name: domain.Name,
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

func toEntity(domain *libvirt.Domain) (Domain, error) {
	domainID, _ := domain.GetUUIDString()
	name, _ := domain.GetName()
	memory, _ := domain.GetMaxMemory() // always in kB
	cpu, _ := domain.GetVcpus()        // temp
	// osType, _ := domain.GetOSType() // temp, ostype doesn't return specific os like ubuntu, centos, etc. just return hvm or other type of vm

	xmlDesc, err := domain.GetXMLDesc(0)
	if err != nil {
		return Domain{}, fmt.Errorf("failed to get XML description: %v", err)
	}

	var domConf libvirtxml.Domain
	if err := domConf.Unmarshal(xmlDesc); err != nil {
		return Domain{}, fmt.Errorf("failed to unmarshal XML description: %v", err)
	}

	return Domain{
		ID:   domainID,
		Name: name,
		Memory: Memory{
			Value: uint(memory / 1024),
			Unit:  UnitMB,
		},
		Cpu: Cpu{
			Value: uint(cpu[0].Cpu), // just for temp
		},

		OS: OS{
			Type: domConf.OS.Type.Type,
			Arch: domConf.OS.Type.Arch,
		},
	}, nil
}
