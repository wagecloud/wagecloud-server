package model


type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type Account interface {
	Base() AccountBase
}

type AccountBase struct {
	ID        int64     `json:"id"` /* unique */
	Role      Role      `json:"role"`
	Name      string    `json:"name"`
	Username  string    `json:"username"` /* unique */
	Password  string    `json:"-"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (a AccountBase) Base() AccountBase {
	return a
}

type AccountUser struct {
	AccountBase
	Email string `json:"email"` /* unique */
}
