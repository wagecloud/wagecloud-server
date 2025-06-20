package libvirt

import (
	"crypto/rand"
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

type Status string

const (
	StatusUnknown Status = "STATUS_UNKNOWN"
	StatusPending Status = "STATUS_PENDING"
	StatusRunning Status = "STATUS_RUNNING"
	StatusStopped Status = "STATUS_STOPPED"
	StatusError   Status = "STATUS_ERROR"
)

func ToStatus(state libvirt.DomainState) Status {
	status := StatusUnknown

	switch state {
	case libvirt.DOMAIN_RUNNING:
		status = StatusRunning
	case libvirt.DOMAIN_SHUTOFF:
		status = StatusStopped
	case libvirt.DOMAIN_SHUTDOWN:
		status = StatusPending // VM is being shut down
	case libvirt.DOMAIN_PAUSED, libvirt.DOMAIN_PMSUSPENDED:
		status = StatusStopped // Treated as "stopped" since VM is not actively running
	case libvirt.DOMAIN_BLOCKED:
		status = StatusPending // VM is waiting on a resource
	case libvirt.DOMAIN_CRASHED:
		status = StatusError // VM crashed unexpectedly
	case libvirt.DOMAIN_NOSTATE:
		status = StatusUnknown
	default:
		status = StatusUnknown
	}

	return status
}

type DomainMonitor struct {
	ID           string
	Status       Status
	CPUUsage     float64 // in percentage
	RAMUsage     float64 // in MB
	StorageUsage float64 // in MB
	NetworkIn    float64 // in MB
	NetworkOut   float64 // in MB
}

type Domain struct {
	ID      string
	Name    string
	Memory  Memory
	Cpu     Cpu
	OS      OS
	Storage uint
	Network DomainNetwork
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

type DomainNetwork struct {
	MacAddress string
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

func GenerateMacAddress() string {
	mac, err := generateRandomMAC()
	if err != nil {
		return ""
	}

	return mac
}

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
		Network: DomainNetwork{
			MacAddress: domainXML.Devices.Interfaces[0].MAC.Address,
		},
	}, nil
}

func getXMLConfig(domain Domain) (*libvirtxml.Domain, error) {
	vmImagePath := domain.VMImagePath()
	cloudinitPath := domain.CloudinitPath()

	if !file.Exists(vmImagePath) || !file.Exists(cloudinitPath) {
		return nil, fmt.Errorf("image or cloudinit file not found")
	}

	var err error
	mac := domain.Network.MacAddress

	if domain.Network.MacAddress == "" {
		mac, err = generateRandomMAC()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random MAC address: %v", err)
		}
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
			MemBalloon: &libvirtxml.DomainMemBalloon{
				Model: "virtio",
				Stats: &libvirtxml.DomainMemBalloonStats{
					Period: 2, // poll every 2 seconds
				},
			},
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
						Address: mac,
					},
					Source: &libvirtxml.DomainInterfaceSource{
						Bridge: &libvirtxml.DomainInterfaceSourceBridge{
							Bridge: "virbr0",
						},
					},
					Model: &libvirtxml.DomainInterfaceModel{
						Type: "virtio",
					},
					Target: &libvirtxml.DomainInterfaceTarget{
						// TODO: try change into vnet1 and see if getPrivateIP works?
						Dev: "vnet0",
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

// generateRandomMAC generates a random MAC address
func generateRandomMAC() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	// Set the locally administered bit (bit 1 of first octet)
	// and clear the multicast bit (bit 0 of first octet)
	buf[0] = (buf[0] | 0x02) & 0xFE

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5]), nil
}
