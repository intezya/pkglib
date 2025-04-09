package main

import (
	"fmt"
	"github.com/intezya/pkglib/crypto"
)

func main() {
	// Example: Hash a value using SHA256
	hashSHA256 := crypto.HashSHA256("hello world")
	fmt.Printf("SHA256 Hash: %s\n", hashSHA256)

	// Example: Hash a value using SHA512
	hashSHA512 := crypto.HashSHA512("hello world")
	fmt.Printf("SHA512 Hash: %s\n", hashSHA512)
}
