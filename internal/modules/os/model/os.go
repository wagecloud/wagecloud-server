package model

import "fmt"

type OS struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

func (o OS) ImageName() string {
	return fmt.Sprintf("%s.img", o.Name)
}
