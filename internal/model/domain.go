package model

import "github.com/google/uuid"

type Memory struct {
	Value uint
	Unit  Unit
}

type Cpu struct {
	Value uint
}

type Domain struct {
	Name          string
	UUID          string
	Arch          Arch
	Memory        Memory
	Cpu           Cpu
	OS            OS
	SourcePath    string
	CloudinitPath string
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

func NewDomain(options ...DomainOption) *Domain {
	// Initialize the domain struct
	domain := &Domain{
		UUID: uuid.New().String(),
	}

	for _, option := range options {
		option(domain)
	}

	return domain
}
