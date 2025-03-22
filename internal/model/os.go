package model

import "fmt"

type OS struct {
	Name string `json:"name"`
	Arch Arch   `json:"arch"`
}

func (o OS) ImageName() string {
	return fmt.Sprintf("%s_%s.img", o.Name, o.Arch)
}
