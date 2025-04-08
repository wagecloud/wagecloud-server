package libvirt

import (
	"fmt"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/qemu"
	"libvirt.org/go/libvirt"
)

var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repo *repository.RepositoryImpl
	qemu *qemu.Service
}

type ServiceInterface interface {
	StartDomain(domain *libvirt.Domain) error
	StartDomainByID(domainID string) error
	CreateDomain(domain model.Domain) (*libvirt.Domain, error)
	UpdateDomain(domainID string, domain model.Domain) (*libvirt.Domain, error)
	GetXMLConfig(domain model.Domain) (*libvirtxml.Domain, error)
	GetDomain(domainID string) (*model.Domain, error)
	GetListDomains() ([]model.Domain, error)
}

func NewService(repo *repository.RepositoryImpl) *Service {
	return &Service{repo: repo}
}

func (s *Service) StartDomainByID(domainID string) error {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	domain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return fmt.Errorf("failed to find domain: %v", err)
	}

	err = domain.Create()

	if err != nil {
		return err
	}

	return nil
}


func (s *Service) StartDomain(domain *libvirt.Domain) error {
	err := domain.Create()

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateDomain(domain model.Domain) (*libvirt.Domain, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	// Create new qcow2 image from base image

	// fake account uuid
	accountID := "7a4a5c55-000c-44d5-b41e-903b71bf32fe"

	if err = s.qemu.CreateImageWithPath(
		domain.BaseImagePath("focal-server-cloudimg-amd64.img"),
		domain.ImageAccountPath(accountID),
		domain.Storage,
	); err != nil {
		return nil, fmt.Errorf("failed to clone image: %v", err)
	}

	domainXML, err := s.GetXMLConfig(domain)

	if err != nil {
		return nil, fmt.Errorf("failed to generate domain XML: %v", err)
	}

	xmlData, err := domainXML.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain XML: %v", err)
	}
	 
	domainVirt, err := conn.DomainDefineXML(xmlData)

	if err != nil {
		return nil, fmt.Errorf("failed to define domain: %v", err)
	}

	return domainVirt, nil

}

func (s *Service) UpdateDomain(domainID string, domain model.Domain) (*libvirt.Domain, error) {

	if domainID == "" {
		return nil, fmt.Errorf("domain ID is required")
	}

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	domainObj, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to find domain: %v", err)
	}

	domainXML, err := s.GetXMLConfig(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to generate domain XML: %v", err)
	}

	xmlData, err := domainXML.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain XML: %v", err)
	}

	newDom, err := conn.DomainDefineXML(xmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to define domain: %v", err)
	}

	// If the domain is running, we need to start it again
	if val, err := domainObj.IsActive(); err != nil && val == true {
		if err := domainObj.Destroy(); err != nil {
			return nil, fmt.Errorf("failed to destroy domain: %v", err)
		}
	}

	return newDom, nil
}

func (s *Service) GetXMLConfig(domain model.Domain) (*libvirtxml.Domain, error) {
	domainXML := &libvirtxml.Domain{
		Type: "kvm",
		Name: domain.Name,
		UUID: domain.UUID,
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
				// Arch:    string(domain.OS.Arch),
				Arch:    "x86_64",
				Machine: "pc-q35-6.2",
				Type:    "hvm",
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
							File: domain.ImageAccountPath("7a4a5c55-000c-44d5-b41e-903b71bf32fe"), // temp accountID := "7a4a5c55-000c-44d5-b41e-903b71bf32fe"
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
							File: domain.CloudinitAccountPath("7a4a5c55-000c-44d5-b41e-903b71bf32fe"), // temp accountID := "7a4a5c55-000c-44d5-b41e-903b71bf32fe"
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

func (s *Service) GetDomain(domainID string) (*model.Domain, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	domain, err := conn.LookupDomainByUUIDString(domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to find domain: %v", err)
	}

	return toEntity(domain)
}

func (s *Service) GetListDomains() ([]model.Domain, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(0)

	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %v", err)
	}

	domainsModel := make([]model.Domain, len(domains))

	for i, domain := range domains {
		model, err := toEntity(&domain)
		if err != nil {
			return nil, fmt.Errorf("failed to convert domain to model: %v", err)
		}
		domainsModel[i] = *model

	}

	return domainsModel, nil
}

func (s *Service) GetListActiveDomains() ([]model.Domain, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %v", err)
	}

	domainsModel := make([]model.Domain, len(domains))

	for i, domain := range domains {
		model, err := toEntity(&domain)
		if err != nil {
			return nil, fmt.Errorf("failed to convert domain to model: %v", err)
		}
		domainsModel[i] = *model
	}

	return domainsModel, nil
}

func toEntity(domain *libvirt.Domain) (*model.Domain, error) {
	domainID, _ := domain.GetUUIDString()
	name, _ := domain.GetName()
	memory, _ := domain.GetMaxMemory() // always in kB
	cpu, _ := domain.GetVcpus()        // temp
	// osType, _ := domain.GetOSType() // temp, ostype doesn't return specific os like ubuntu, centos, etc. just return hvm or other type of vm

	xmlDesc, err := domain.GetXMLDesc(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get XML description: %v", err)
	}

	var domConf libvirtxml.Domain
	if err := domConf.Unmarshal(xmlDesc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML description: %v", err)
	}

	return &model.Domain{
		UUID: domainID,
		Name: name,
		Memory: model.Memory{
			Value: uint(memory / 1024),
			Unit:  model.UnitMB,
		},
		Cpu: model.Cpu{
			Value: uint(cpu[0].Cpu), // just for temp
		},

		OS: model.OS{
			Name: domConf.OS.Type.Type,
			// Arch: model.Arch(domConf.OS.Type.Arch),
		}}, nil

}
