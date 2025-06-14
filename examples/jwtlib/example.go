package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"

	"github.com/intezya/pkglib/jwtlib"
)

type MyClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
	UID  int64  `json:"uid"`
}

func main() {
	tm := jwtlib.New[MyClaims](
		jwtlib.Config{
			SecretKey:     "secret",
			Issuer:        "issuer",
			TokenDuration: time.Hour,
			Strict:        false,
		},
		func() MyClaims { return MyClaims{} },
	)

	claims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user@example.com",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
		Role: "admin",
		UID:  42,
	}

	tokenStr := tm.Generate(claims)
	fmt.Println("Token:", tokenStr)

	parsed, err := tm.Parse(tokenStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsed:", parsed.Subject, parsed.UID, parsed.Role)
}
