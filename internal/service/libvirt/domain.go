package libvirt

import (
	"fmt"
	"path"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/model"
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
	return fmt.Sprintf("vm_%s.iso", d.ID)
}

func (d Domain) BaseFileName() string {
	return fmt.Sprintf("%s_%s.iso", d.OS.Arch, d.OS.Name)
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

func ToDomain(vm model.VM) Domain {
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
