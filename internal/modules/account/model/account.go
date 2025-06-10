package accountmodel

import "github.com/golang-jwt/jwt/v5"

type Role string

const (
	RoleAdmin Role = "ROLE_ADMIN"
	RoleStaff Role = "ROLE_STAFF"
	RoleUser  Role = "ROLE_USER"
)

type Account interface {
	Base() AccountBase
}

type AccountBase struct {
	ID        int64  `json:"id"` /* unique */
	Role      Role   `json:"role"`
	Name      string `json:"name"`
	Username  string `json:"username"` /* unique */
	Password  string `json:"-"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (a AccountBase) Base() AccountBase {
	return a
}

type AccountUser struct {
	AccountBase
	Email string `json:"email"` /* unique */
}

type AuthenticatedAccount struct {
	AccountID int64 `json:"account_id"`
	Role      Role  `json:"role"`
}

type Claims struct {
	AccountID int64
	Role      Role
	jwt.RegisteredClaims
}

func (c *Claims) ToAuthenticatedAccount() AuthenticatedAccount {
	return AuthenticatedAccount{
		AccountID: c.AccountID,
		Role:      c.Role,
	}
}
