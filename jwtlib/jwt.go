package jwtlib

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenManager[T Claims] struct {
	secretKey     []byte
	issuer        string
	tokenDuration time.Duration
	strict        bool
	newClaims     func() T
}

// New creates a new token manager with custom claims
func New[T Claims](config Config, newClaims func() T) *TokenManager[T] {
	if len(config.SecretKey) < 32 {
		panic("secret key must be at least 32 bytes")
	}
	if newClaims == nil {
		panic("newClaims factory required")
	}

	return &TokenManager[T]{
		secretKey:     []byte(config.SecretKey),
		issuer:        config.Issuer,
		tokenDuration: config.TokenDuration,
		strict:        config.Strict,
		newClaims:     newClaims,
	}
}

// Generate creates a signed token using the provided claims
func (tm *TokenManager[T]) Generate(claims T) string {
	now := time.Now()
	claims.SetIssuedAt(jwt.NewNumericDate(now))
	claims.SetExpiresAt(jwt.NewNumericDate(now.Add(tm.tokenDuration)))
	claims.SetIssuer(tm.issuer)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(tm.secretKey)

	if err != nil {
		panic("unexpected jwt signing error: " + err.Error())
	}

	return signed
}

// Parse validates token and returns custom claims
func (tm *TokenManager[T]) Parse(tokenStr string) (T, error) {
	var parserOpts []jwt.ParserOption

	if tm.strict {
		parserOpts = append(parserOpts, jwt.WithStrictDecoding())
	}

	parser := jwt.NewParser(parserOpts...)
	claims := tm.newClaims()

	token, err := parser.ParseWithClaims(
		tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return tm.secretKey, nil
		},
	)
	if err != nil {
		var zero T
		return zero, err
	}

	if !token.Valid {
		var zero T
		return zero, errors.New("invalid token")
	}

	return claims, nil
}
