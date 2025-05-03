package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

// Salt generates a random salt of the given length.
// It returns the generated salt and an error if it fails.
func Salt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// HashSHA256 hashes the input string using SHA-256 and returns the hexadecimal representation.
func HashSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

// HashSHA512 hashes the input string using SHA-512 and returns the hexadecimal representation.
func HashSHA512(value string) string {
	hash := sha512.Sum512([]byte(value))
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

// EncodeBase64 encodes the input string into a base64-encoded string.
func EncodeBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// DecodeBase64 decodes a base64-encoded string into its original string representation.
// If the decoding fails, it returns an error.
func DecodeBase64(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	return string(data), nil
}

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultArgonParams = &ArgonParams{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	return salt, err
}

func HashArgon2(password string, p *ArgonParams) (string, error) {
	salt, err := generateSalt(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory,
		p.Iterations,
		p.Parallelism,
		saltB64,
		hashB64,
	)

	return encoded, nil
}

func VerifyArgon2(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	calculated := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))
	return subtleCompare(hash, calculated), nil
}

func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
