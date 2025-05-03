package generate

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"math/big"
)

// RandomString generates a random string of the given length using the specified charset.
// The charset string defines the allowed characters to form the random string.
func RandomString(length int, charset string) string {
	result := make([]byte, length)
	mx := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, mx)
		if err != nil {
			panic("failed to generate random string: " + err.Error())
		}
		result[i] = charset[n.Int64()]
	}

	return string(result)
}

// RandomSHA256 generates a random string and returns its SHA-256 hash in hexadecimal form.
func RandomSHA256() string {
	hash := sha256.Sum256([]byte(RandomString(32, AllCharset)))
	return hex.EncodeToString(hash[:])
}

// RandomSHA512 generates a random string and returns its SHA-512 hash in hexadecimal form.
func RandomSHA512() string {
	hash := sha512.Sum512([]byte(RandomString(32, AllCharset)))
	return hex.EncodeToString(hash[:])
}

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic("failed to generate random bytes: " + err.Error())
	}
	return b
}
