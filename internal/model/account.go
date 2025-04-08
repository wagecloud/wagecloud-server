package model

import "time"

type Account interface {
	Base() AccountBase
}

type AccountBase struct {
	ID        int64     `json:"id"`       /* unique */
	Username  string    `json:"username"` /* unique */
	Email     string    `json:"email"`    /* unique */
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (a AccountBase) Base() AccountBase {
	return a
}

type AccountUser struct {
	AccountBase
}
