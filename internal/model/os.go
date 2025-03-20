package model

type OSType string

const (
	OSTypeUbuntu OSType = "ubuntu"
	OSTypeDebian OSType = "debian"
)

type OS struct {
	Arch Arch
	Type OSType
}
