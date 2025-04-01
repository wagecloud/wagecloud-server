package model

import "time"

type AccountBase struct {
	ID        int64     `json:"id"`    /* unique */
	Email     string    `json:"email"` /* unique */
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
