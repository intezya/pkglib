package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Service represents a JWT service that provides functionality for
// validating and generating JWT tokens. It contains the secret key,
// issuer, and expiration time used during token operations.
type Service struct {
	Configuration
}

// New creates and returns a new instance of Service,
// initialized with values from the provided Configuration.
func New(config Configuration) *Service {
	return &Service{Configuration: config}
}

// Validate parses and validates the given JWT token string.
// It returns the authentication data if the token is valid, or an error if invalid.
func (s *Service) Validate(tokenString string, customOptions ...jwt.ParserOption) (
	authData AuthenticationData,
	err error,
) {
	claims := &Claim{}

	if len(customOptions) == 0 {
		customOptions = []jwt.ParserOption{
			jwt.WithIssuer(s.issuer), // Validate the token's issuer
		}
	}

	// Parse the token and extract claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for validating the token signature
			return s.secretKey, nil
		},
		customOptions...,
	)

	// If there's an error parsing or validating the token, return an error
	if err != nil {
		return nil, err
	}

	// Ensure the claims are of the expected type and that the token is valid
	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the authentication data embedded in the token's claims
	return claims.AuthenticationData, nil
}

// Generate creates a new JWT token string with the provided authentication data.
// It returns the generated token string and any error encountered during creation.
func (s *Service) Generate(authData AuthenticationData) (tokenString string, err error) {
	// Define the claims, including authentication data and registered claims (e.g., expiration)
	claims := &Claim{
		AuthenticationData: authData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expirationTime)), // Set token expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       // Set token issue time
			NotBefore: jwt.NewNumericDate(time.Now()),                       // Set token not-before time
			Issuer:    s.issuer,                                             // Set token issuer
		},
	}

	// Create the token with the claims and the signing method (HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Return the signed token string using the secret key
	return token.SignedString(s.secretKey)
}
