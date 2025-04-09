package jwt

import "github.com/golang-jwt/jwt/v5"

// AuthenticationData is a type alias for a map that holds key-value pairs of string data
// that represent the authentication information encoded in the JWT.
type AuthenticationData = map[string]string

// Claim represents the structure of a JWT claim. It includes both the authentication data and
// the registered claims (e.g., expiration, issuer).
// RegisteredClaims is embedded from the golang-jwt package.
type Claim struct {
	AuthenticationData `json:"authentication_data"`
	jwt.RegisteredClaims
}
