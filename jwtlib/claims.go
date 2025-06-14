package jwtlib

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims interface {
	jwt.Claims
	SetIssuer(string)
	SetExpiresAt(*jwt.NumericDate)
	SetIssuedAt(*jwt.NumericDate)
}
