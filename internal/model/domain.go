package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/wagecloud/wagecloud-server/config"
)

type Memory struct {
	Value uint
	Unit  Unit
}

type Cpu struct {
	Value uint
}

type Domain struct {
	UUID   string
	Name   string
	Memory Memory
	Cpu    Cpu
	OS     OS
}

func (d Domain) BaseImagePath() string {
	return fmt.Sprintf("%s/%s", config.GetConfig().App.BaseImageDir, d.OS.ImageName())
}

func (d Domain) ImagePath() string {
	return fmt.Sprintf("%s/%s.img", config.GetConfig().App.ImageDir, d.UUID)
}

func (d Domain) CloudinitPath() string {
	return fmt.Sprintf("%s/cloudinit_%s.iso", config.GetConfig().App.CloudinitDir, d.UUID)
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
