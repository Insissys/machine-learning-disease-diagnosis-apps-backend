package model

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	Email   string
	Role    string
	GroupID uint64
	Revoked bool

	jwt.RegisteredClaims
}
