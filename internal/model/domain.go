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
	UUID   string
	Name   string
	Memory Memory
	Cpu    Cpu
	OS     OS
	Storage uint 
}

func (d Domain) BaseImagePath() string {
	return path.Join(
		config.GetConfig().App.BaseImageDir,
		d.OS.ImageName(),
	)
}

func (d Domain) ImagePath() string {
	return path.Join(
		config.GetConfig().App.ImageDir,
		fmt.Sprintf("%s.img", d.UUID),
	)
}

func (d Domain) CloudinitPath() string {
	return path.Join(
		config.GetConfig().App.CloudinitDir,
		fmt.Sprintf("cloudinit_%s.iso", d.UUID),
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
