package libvirt

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"github.com/wagecloud/wagecloud-server/internal/service/qemu"
)

type Service struct {
	repo *repository.Repository
	qemu *qemu.Service
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
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

	if err = s.qemu.CreateImage(domain.BaseImagePath(), domain.ImagePath(), domain.Storage); err != nil {
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

	return conn.DomainDefineXML(xmlData)
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
				Arch:    string(domain.OS.Arch),
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
							File: domain.ImagePath(),
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
							File: domain.CloudinitPath(),
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
