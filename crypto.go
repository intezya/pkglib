package pkglib

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/argon2"
)

type crypto struct {
}

func (crypto) Salt(length int) []byte {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	return salt
}

func (crypto) HashSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

func (crypto) HashSHA512(value string) string {
	hash := sha512.Sum512([]byte(value))
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

func (crypto) EncodeBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func (crypto) DecodeBase64(value string) string {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return ""
	}
	return string(data)
}

func (crypto) HashArgon2(password string, salt []byte) string {
	timeCost := uint32(1)           // Number of iterations
	memoryCost := uint32(64 * 1024) // Memory in KB (64 MB)
	parallelism := uint8(2)         // Number of threads
	keyLength := uint32(32)         // Length of the generated hash

	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, keyLength)

	return fmt.Sprintf(
		"%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

func (crypto) VerifyArgon2(password, encodedHash string) bool {
	parts := bytes.Split([]byte(encodedHash), []byte("$"))
	if len(parts) != 2 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(string(parts[0]))
	if err != nil {
		return false
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(string(parts[1]))
	if err != nil {
		return false
	}

	timeCost := uint32(1)
	memoryCost := uint32(64 * 1024)
	parallelism := uint8(2)
	keyLength := uint32(32)

	actualHash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, keyLength)

	return bytes.Equal(actualHash, expectedHash)
}

var Crypto crypto
