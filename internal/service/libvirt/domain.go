package libvirt

import (
	"fmt"
	"path"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/model"
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

func FromVMToDomain(vm model.VM) Domain {
	return Domain{
		ID:   vm.ID,
		Name: vm.Name,
		Memory: Memory{
			Value: uint(vm.Ram),
			Unit:  "MB",
		},
		Cpu: Cpu{
			Value: uint(vm.Cpu),
		},
		OS: OS{
			Name: vm.OsID,
			Type: "hvm",
			Arch: vm.ArchID,
		},
		Storage: uint(vm.Storage),
	}
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
	}, nil
}
