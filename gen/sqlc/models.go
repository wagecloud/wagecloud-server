// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentType string

const (
	PaymentTypeVNPay PaymentType = "VNPay"
)

func (e *PaymentType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentType(s)
	case string:
		*e = PaymentType(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentType: %T", src)
	}
	return nil
}

type NullPaymentType struct {
	PaymentType PaymentType
	Valid       bool // Valid is true if PaymentType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentType) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentType), nil
}

type Role string

const (
	RoleADMIN Role = "ADMIN"
	RoleUSER  Role = "USER"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role
	Valid bool // Valid is true if Role is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type AccountBase struct {
	ID        int64
	Role      Role
	Name      string
	Username  string
	Password  string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type AccountUser struct {
	ID    int64
	Email string
}

type Arch struct {
	ID        string
	Name      string
	CreatedAt pgtype.Timestamptz
}

type Network struct {
	ID        string
	PrivateIp string
	CreatedAt pgtype.Timestamptz
}

type O struct {
	ID        string
	Name      string
	CreatedAt pgtype.Timestamptz
}

type Vm struct {
	ID        string
	AccountID int64
	NetworkID string
	OsID      string
	ArchID    string
	Name      string
	Cpu       int32
	Ram       int32
	Storage   int32
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}
