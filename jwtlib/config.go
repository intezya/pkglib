package jwtlib

import "time"

type Config struct {
	SecretKey     string        `env:"JWT_SECRET_KEY"`
	Issuer        string        `env:"JWT_ISSUER"`
	TokenDuration time.Duration `env:"JWT_TOKEN_DURATION" env-default:"1d"`
	Strict        bool          `env:"JWT_STRICT_DECODING" env-default:"false"`
}
