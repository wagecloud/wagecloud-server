package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	AccountID int64
	jwt.RegisteredClaims
}
