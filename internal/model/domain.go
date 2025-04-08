package model

import (
	"fmt"
	"path"

	"github.com/google/uuid"
	"github.com/wagecloud/wagecloud-server/config"
)

type Memory struct {
	Value uint `json:"value"`
	Unit  Unit `json:"unit"`
}

type Cpu struct {
	Value uint `json:"value"`
}

type Domain struct {
	UUID    string
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
		fmt.Sprintf("%s.img", d.UUID),
	)
}

func (d Domain) ImageAccountPath(accountID string) string {
	return path.Join(
		config.GetConfig().App.ImageDir,
		accountID,
		fmt.Sprintf("%s.img", d.UUID),
	)
}

func (d Domain) CloudinitPath() string {
	return path.Join(
		config.GetConfig().App.CloudinitDir,
		fmt.Sprintf("cloudinit_%s.iso", d.UUID),
	)
}

func (d Domain) CloudinitAccountPath(accountID string) string {
	return path.Join(
		config.GetConfig().App.CloudinitDir,
		accountID,
		"ubuntu-with-init.iso",
	)
}

type DomainOption func(*Domain)

func WithDomainName(name string) DomainOption {
	return func(domain *Domain) {
		domain.Name = name
	}
}

func WithDomainMemory(value uint, unit Unit) DomainOption {
	return func(domain *Domain) {
		domain.Memory = Memory{Value: value, Unit: unit}
	}
}

func WithDomainCpu(value uint) DomainOption {
	return func(domain *Domain) {
		domain.Cpu = Cpu{Value: value}
	}
}

func WithDomainOS(os OS) DomainOption {
	return func(domain *Domain) {
		domain.OS = os
	}
}

func WithDomainStorage(storage uint) DomainOption {
	return func(domain *Domain) {
		domain.Storage = storage
	}
}

func NewDomain(options ...DomainOption) Domain {
	// Initialize the domain struct
	domain := &Domain{
		UUID: uuid.New().String(),
	}

	for _, option := range options {
		option(domain)
	}

	return *domain
}
