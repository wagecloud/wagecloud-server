package core

import (
	"GoVirService/internal/model"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

type InstallDomainParams struct {
	OS            model.OS
	Memory        int
	VCPUs         uint
	SourcePath    string
	CloudInitPath string
}

func InstallDomain(params InstallDomainParams) {
	// Connect to libvirt
	conn, err := libvirt.NewConnect("qemu:///system")

	if err != nil {
		log.Fatalf("Failed to connect to libvirt: %v", err)
	}

	defer conn.Close()

	fmt.Println("Connected to libvirt")

	id := uuid.New()

	domainXML := &libvirtxml.Domain{
		Type: "kvm",
		// Name: fmt.Sprintf("%s-%s", params.OS.Name, id.String()),
		Name: "ubuntu",
		UUID: id.String(),
		Memory: &libvirtxml.DomainMemory{
			Value: 2048,
			Unit:  "MiB",
		},

		// If this is omitted, it defaults to the same value as the memory element. The unit attribute behaves the same as for memory.
		// CurrentMemory: &libvirtxml.DomainCurrentMemory{
		// 	Value: 2048,
		// 	Unit:  "MiB",
		// },
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			Value:     params.VCPUs,
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    "x86_64",
				Machine: "pc",
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
						Name: "qemu",
						Type: "raw", // Use raw instead of qcow2
					},
					Source: &libvirtxml.DomainDiskSource{
						File: &libvirtxml.DomainDiskSourceFile{
							File: params.SourcePath,
						},
						Index: 1,
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
							File: params.CloudInitPath,
						},
						Index: 1,
					},
					Target: &libvirtxml.DomainDiskTarget{
						Dev: "sda",
						Bus: "sata",
					},
					ReadOnly: &libvirtxml.DomainDiskReadOnly{},
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
					// MAC: &libvirtxml.DomainInterfaceMAC{
					// 	Address: "52:54:00:b7:a5:c2",
					// },

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

	if err := domain.Create(); err != nil {
		log.Fatalf("Failed to start VM: %v", err)
	}

	fmt.Println("Domain started")
}
