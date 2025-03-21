package model

import "fmt"

type OS struct {
	Name string // e.g. Ubuntu, CentOS, Debian
	Arch Arch
}

func (o OS) ImageName() string {
	return fmt.Sprintf("%s_%s.img", o.Name, o.Arch)
}
