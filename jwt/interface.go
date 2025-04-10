package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Validate defines the interface for validating a JWT token string.
// It returns the authentication data and an error if the token is invalid.
type Validate interface {
	Validate(tokenString string, customOptions ...jwt.ParserOption) (authData AuthenticationData, err error)
}

// Generate defines the interface for generating a new JWT token string
// based on the provided authentication data.
type Generate interface {
	Generate(authData AuthenticationData) (tokenString string, err error)
}

// Configuration is a struct for providing configuration values
// needed for JWT generation and validation (secret key, issuer, expiration time).
type Configuration struct {
	secretKey      []byte
	issuer         string
	expirationTime time.Duration
}

func NewConfiguration(secretKey []byte, issuer string, expirationTime time.Duration) Configuration {
	return Configuration{
		secretKey:      secretKey,
		issuer:         issuer,
		expirationTime: expirationTime,
	}
}

// TokenService combines the Validate and Generate interfaces.
// It represents a service that can both validate and generate JWT tokens.
type TokenService interface {
	Validate
	Generate
}
