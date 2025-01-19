package pkglib

import (
	"errors"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

type jwt struct {
}

func (jwt) GenerateToken(payload map[string]interface{}, secret string) string {
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.MapClaims(payload))
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return signedToken
}

func (jwt) ValidateJWT(token string, secret string) (map[string]interface{}, error) {
	claims := jwtlib.MapClaims{}
	_, err := jwtlib.ParseWithClaims(
		token, claims, func(t *jwtlib.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]interface{})
	for key, val := range claims {
		result[key] = val
	}
	return result, nil
}

var JWT jwt
