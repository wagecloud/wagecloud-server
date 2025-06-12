package accountmodel

import "github.com/golang-jwt/jwt/v5"

type AccountType string

const (
	AccountTypeAdmin AccountType = "ACCOUNT_TYPE_ADMIN"
	AccountTypeUser  AccountType = "ACCOUNT_TYPE_USER"
)

type Account interface {
	Base() AccountBase
}

type AccountBase struct {
	ID        int64       `json:"id"` /* unique */
	Type      AccountType `json:"type"`
	Name      string      `json:"name"`
	Username  string      `json:"username"` /* unique */
	Password  string      `json:"-"`
	CreatedAt int64       `json:"created_at"`
	UpdatedAt int64       `json:"updated_at"`
}

func (a AccountBase) Base() AccountBase {
	return a
}

type AccountUser struct {
	AccountBase
	Email string `json:"email"` /* unique */
}

type AuthenticatedAccount struct {
	AccountID int64       `json:"account_id"`
	Type      AccountType `json:"type"`
}

type Claims struct {
	AccountID int64
	Type      AccountType
	jwt.RegisteredClaims
}

func (c *Claims) ToAuthenticatedAccount() AuthenticatedAccount {
	return AuthenticatedAccount{
		AccountID: c.AccountID,
		Type:      c.Type,
	}
}
