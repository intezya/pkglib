package crypto

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

// ArgonConfig holds the configuration settings for Argon2 hashing.
// It includes parameters like time cost, memory cost, parallelism, and key length.
type ArgonConfig struct {
	TimeCost    uint32
	MemoryCost  uint32
	Parallelism uint8
	KeyLength   uint32
}

// defaultArgonConfig provides a default configuration for Argon2 hashing.
var defaultArgonConfig = &ArgonConfig{
	TimeCost:    1,
	MemoryCost:  64 * 1024,
	Parallelism: 2,
	KeyLength:   32,
}

// HashArgon2 hashes a password using the Argon2 ID variant.
// It returns the salted hash as a base64-encoded string, formatted as "<salt>$<hash>".
func HashArgon2(password string, salt []byte, config *ArgonConfig) string {
	if config == nil {
		config = defaultArgonConfig
	}

	// Perform Argon2 hashing
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		config.TimeCost,
		config.MemoryCost,
		config.Parallelism,
		config.KeyLength,
	)

	// Return the result as "<salt>$<hash>"
	return fmt.Sprintf(
		"%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

// VerifyArgon2 verifies a password against an Argon2 hash.
// It returns true if the password matches the hash, otherwise false.
func VerifyArgon2(password, encodedHash string, config *ArgonConfig) bool {
	if config == nil {
		config = defaultArgonConfig
	}

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

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		config.TimeCost,
		config.MemoryCost,
		config.Parallelism,
		config.KeyLength,
	)

	return bytes.Equal(actualHash, expectedHash)
}
