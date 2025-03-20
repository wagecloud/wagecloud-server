package model

type Memory struct {
	Value int
	Unit  Unit
}

type Cpu struct {
	Value int
}

type Domain struct {
	Name   string
	UUID   string
	Memory Memory
	Cpu    Cpu
	OS     OS
}
