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
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(charsetLen)]
	}
	return string(code)
}

// RandomSHA256 generates a random string and returns its SHA-256 hash in hexadecimal form.
func RandomSHA256() string {
	hash := sha256.Sum256([]byte(RandomString(32, Charset.AllCharset)))
	return hex.EncodeToString(hash[:])
}

// RandomSHA512 generates a random string and returns its SHA-512 hash in hexadecimal form.
func RandomSHA512() string {
	hash := sha512.Sum512([]byte(RandomString(32, Charset.AllCharset)))
	return hex.EncodeToString(hash[:])
}
