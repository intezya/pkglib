package generate

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"time"
)

// RandomString generates a random string of the given length using the specified charset.
// The charset string defines the allowed characters to form the random string.
func RandomString(length int, charset string) string {
	charsetLen := len(charset)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed random generator with current timestamp
	code := make([]byte, length)
	// Generate random characters from the charset
	for i := range code {
		code[i] = charset[seededRand.Intn(charsetLen)] // Pick a random character from charset
	}

	return string(code)
}

// RandomSHA256 generates a random string and returns its SHA-256 hash in hexadecimal form.
func RandomSHA256() string {
	// Generate a random string of length 32 and hash it using SHA-256
	hash := sha256.Sum256([]byte(RandomString(32, Charset.AllCharset)))
	return hex.EncodeToString(hash[:]) // Convert the hash to a hexadecimal string
}

// RandomSHA512 generates a random string and returns its SHA-512 hash in hexadecimal form.
func RandomSHA512() string {
	// Generate a random string of length 32 and hash it using SHA-512
	hash := sha512.Sum512([]byte(RandomString(32, Charset.AllCharset)))
	return hex.EncodeToString(hash[:]) // Convert the hash to a hexadecimal string
}
