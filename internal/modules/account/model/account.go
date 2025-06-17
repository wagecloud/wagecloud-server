package accountmodel

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccountType string

const (
	AccountTypeAdmin AccountType = "ACCOUNT_TYPE_ADMIN"
	AccountTypeUser  AccountType = "ACCOUNT_TYPE_USER"
)

type AccountBase struct {
	ID        int64       `json:"id"` /* unique */
	Type      AccountType `json:"type"`
	Username  string      `json:"username"` /* unique */
	Password  string      `json:"-"`
	CreatedAt time.Time   `json:"created_at"`
}

type AccountUser struct {
	ID        int64   `json:"id"` /* unique */
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     *string `json:"email"` /* unique */
	Phone     *string `json:"phone"` /* unique */
	Company   *string `json:"company"`
	Address   *string `json:"address"`
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
