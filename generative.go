package pkglib

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"time"
)

type generative struct {
}

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (generative) RandomString(length int, charset ...string) string {
	var chars string
	if len(charset) > 0 {
		chars = charset[0]
	} else {
		chars = CHARSET
	}
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[seededRand.Intn(len(chars))]
	}
	return string(code)
}

func (g generative) RandomSHA256() string {
	hash := sha256.Sum256([]byte(g.RandomString(32)))
	return hex.EncodeToString(hash[:])
}

func (g generative) RandomSHA512() string {
	hash := sha512.Sum512([]byte(g.RandomString(32)))
	return hex.EncodeToString(hash[:])
}

var Generative generative
