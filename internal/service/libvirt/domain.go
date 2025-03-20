package domain

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

type MemoryUnit string

const (
	MiB MemoryUnit = "MiB"
	GiB MemoryUnit = "GiB"
)

type OSType string

const (
	OSTypeUbuntu OSType = "ubuntu"
	OSTypeDebian OSType = "debian"
)

type Memory struct {
	Value uint
	Unit  MemoryUnit
}

type Spec struct {
	CPU    uint
	Memory *Memory
	OS     OSType
}

func CreateDomain(spec *Spec) (*libvirt.Domain, error) {
	conn, err := libvirt.NewConnect("qemu:///system")

	if err != nil {
		log.Fatalf("Failed to connect to libvirt: %v", err)
	}

	defer conn.Close()

	fmt.Println("Connected to libvirt")

	domainXML := &libvirtxml.Domain{
		Type: "kvm",
		Name: "ubuntu",
		UUID: uuid.New().String(),
		Memory: &libvirtxml.DomainMemory{
			Value: 2048,
			Unit:  "MiB",
		},
		CurrentMemory: &libvirtxml.DomainCurrentMemory{
			Value: 2048,
			Unit:  "MiB",
		},
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			Value:     1,
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    "x86_64",
				Machine: "pc-q35-6.2",
				Type:    "hvm",
				// ID:      "ubuntu20.04",
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
							File: "/home/khoakomlem/wage-cloud/wagecloud-server/testcloudinit/focal-server-cloudimg-amd64.img",
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
							File: "/var/lib/libvirt/images/vm1.iso",
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

	xmlData, err := domainXML.Marshal()
	if err != nil {
		log.Fatalf("Failed to generate XML: %v", err)
	}
	domain, err := conn.DomainDefineXML(xmlData)

	if err != nil {
		log.Fatalf("Failed to define domain: %v", err)
	}

	return domain, nil
}

func StartDomain(domain *libvirt.Domain) error {
	err := domain.Create()

	if err != nil {
		return err
	}

	return nil
}
