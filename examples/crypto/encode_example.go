package main

import (
	"fmt"
	"github.com/intezya/pkglib/crypto"
	"log"
)

func main() {
	// Example: Base64 encoding
	encoded := crypto.EncodeBase64("hello world")
	fmt.Printf("Base64 Encoded: %s\n", encoded)

	// Example: Base64 decoding
	decoded, err := crypto.DecodeBase64(encoded)
	if err != nil {
		log.Fatalf("Error decoding base64: %v", err)
	}
	fmt.Printf("Base64 Decoded: %s\n", decoded)
}
