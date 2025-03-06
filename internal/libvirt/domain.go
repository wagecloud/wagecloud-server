package libvirt

import (
	"fmt"
	"log"
	// "time"

	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
	"github.com/libvirt/libvirt-go-xml"
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
		Name: "debian12",
		UUID: uuid.New().String(),
		Memory: &libvirtxml.DomainMemory{
			Value: spec.Memory.Value,
			Unit:  string(spec.Memory.Unit),
		},
		CurrentMemory: &libvirtxml.DomainCurrentMemory{
			Value: spec.Memory.Value,
			Unit:  string(spec.Memory.Unit),
		},
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			Value:     spec.CPU,
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    "x86_64",
				Machine: "pc-q35-6.2",
				Type:    "hvm",
			},
			BootDevices: []libvirtxml.DomainBootDevice{
				{Dev: "cdrom"}, // Boot from ISO first
				{Dev: "hd"},    // Then hard disk
			},
		},
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
			// APIC:   &libvirtxml.DomainFeature{},
			// VMPort: &libvirtxml.DomainFeatureVM{State: "off"},
		},
		CPU: &libvirtxml.DomainCPU{
			Mode:       "host-passthrough",
			Check:      "none",
			Migratable: "on",
		},
		Clock: &libvirtxml.DomainClock{
			Offset: "utc",
			// Timers: []*libvirtxml.DomainTimer{
			// 	{Name: "rtc", TickPolicy: "catchup"},
			// 	{Name: "pit", TickPolicy: "delay"},
			// 	{Name: "hpet", Present: "no"},
			// },
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
							File: "/var/lib/libvirt/images/debian-12-genericcloud-amd64.qcow2",
						},
						Index: 1,
					},
					Target: &libvirtxml.DomainDiskTarget{
						Dev: "vda",
						Bus: "virtio",
					},
				},
				// {
				// 	Device: "cdrom",
				// 	Driver: &libvirtxml.DomainDiskDriver{
				// 		Name: "qemu",
				// 		Type: "raw",
				// 	},
				// 	Source: &libvirtxml.DomainDiskSource{
				// 		File: &libvirtxml.DomainDiskSourceFile{
				// 			File: "/var/lib/libvirt/images/debian-12.9.0-amd64-netinst.iso",
				// 		},
				// 		Index: 1,
				// 	},
				// 	Target: &libvirtxml.DomainDiskTarget{
				// 		Dev: "sda",
				// 		Bus: "sata",
				// 	},
				// 	ReadOnly: &libvirtxml.DomainDiskReadOnly{},
				// },
			},
			Interfaces: []libvirtxml.DomainInterface{
				{
					MAC: &libvirtxml.DomainInterfaceMAC{
						Address: "52:54:00:b7:a5:c2",
					},

					Source: &libvirtxml.DomainInterfaceSource{
						Network: &libvirtxml.DomainInterfaceSourceNetwork{
							Network: "default",
							Bridge:  "virbr0",
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

					Spice: &libvirtxml.DomainGraphicSpice{
						Port:   -1,
						Listen: "0.0.0.0",
					},
				},
			},
		},
	}

	xmlData, err := domainXML.Marshal()
	if err != nil {
		log.Fatalf("Failed to generate XML: %v", err)
	}

	fmt.Println(string(xmlData))

	domain, err := conn.DomainDefineXML(xmlData)
	if err != nil {
		log.Fatalf("Failed to define domain: %v", err)
	}

	defer domain.Free()

	return domain, nil
}
