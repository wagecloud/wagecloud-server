package libvirt

import (
	"fmt"
	"path"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/utils/file"
	"libvirt.org/go/libvirt"
)

type Memory struct {
	Value uint `json:"value"`
	Unit  Unit `json:"unit"`
}

type Cpu struct {
	Value uint `json:"value"`
}

type OS struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Arch string `json:"arch"`
}

type Domain struct {
	ID      string
	Name    string
	Memory  Memory
	Cpu     Cpu
	OS      OS
	Storage uint
}

func (d Domain) CloudinitFileName() string {
	return fmt.Sprintf("cloudinit_%s.iso", d.ID)
}

func (d Domain) VMFileName() string {
	return fmt.Sprintf("vm_%s.img", d.ID)
}

func (d Domain) BaseFileName() string {
	return fmt.Sprintf("%s_%s.img", d.OS.Name, d.OS.Arch)
}

func (d Domain) CloudinitPath() string {
	return path.Join(config.GetConfig().App.CloudinitDir, d.CloudinitFileName())
}

func (d Domain) VMImagePath() string {
	return path.Join(config.GetConfig().App.VMImageDir, d.VMFileName())
}

func (d Domain) BaseImagePath() string {
	return path.Join(config.GetConfig().App.BaseImageDir, d.BaseFileName())
}

// func FromVMToDomain(vm model.VM) Domain {
// 	return Domain{
// 		ID:   vm.ID,
// 		Name: vm.Name,
// 		Memory: Memory{
// 			Value: uint(vm.Ram),
// 			Unit:  "MB",
// 		},
// 		Cpu: Cpu{
// 			Value: uint(vm.Cpu),
// 		},
// 		OS: OS{
// 			Name: vm.OsID,
// 			Type: "hvm",
// 			Arch: vm.ArchID,
// 		},
// 		Storage: uint(vm.Storage),
// 	}
// }

func FromLibvirtToDomain(domain libvirt.Domain) (Domain, error) {
	domainID, _ := domain.GetUUIDString()
	name, _ := domain.GetName()
	memory, _ := domain.GetMaxMemory() // return kB, should convert to MB later
	vcpus, _ := domain.GetMaxVcpus()
	osType, _ := domain.GetOSType()

	xmlDesc, err := domain.GetXMLDesc(0)
	if err != nil {
		return Domain{}, fmt.Errorf("failed to get XML description: %v", err)
	}

	var domainXML libvirtxml.Domain
	if err := domainXML.Unmarshal(xmlDesc); err != nil {
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
			Value: vcpus,
		},
		OS: OS{
			Type: osType,
			Arch: domainXML.OS.Type.Arch,
		},
	}, nil
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
