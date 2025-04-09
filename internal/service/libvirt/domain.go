package libvirt

import (
	"fmt"
	"path"

	"github.com/wagecloud/wagecloud-server/config"
)

type Memory struct {
	Value uint `json:"value"`
	Unit  Unit `json:"unit"`
}

type Cpu struct {
	Value uint `json:"value"`
}

type OS struct {
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

func (d Domain) BaseImagePath(baseOsFileName string) string {
	return path.Join(
		config.GetConfig().App.BaseImageDir,
		baseOsFileName,
	)
}

func (d Domain) ImagePath() string {
	return path.Join(
		config.GetConfig().App.ImageDir,
		fmt.Sprintf("%s.img", d.ID),
	)
}

func (d Domain) ImageAccountPath(accountID string) string {
	return path.Join(
		config.GetConfig().App.ImageDir,
		accountID,
		fmt.Sprintf("%s.img", d.ID),
	)
}

func (d Domain) CloudinitPath() string {
	return path.Join(
		config.GetConfig().App.CloudinitDir,
		fmt.Sprintf("cloudinit_%s.iso", d.ID),
	)
}

func (d Domain) CloudinitAccountPath(accountID string) string {
	return path.Join(
		config.GetConfig().App.CloudinitDir,
		accountID,
		"ubuntu-with-init.iso",
	)
}
