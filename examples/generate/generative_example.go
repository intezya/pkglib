package main

import (
	"fmt"
	"github.com/intezya/pkglib/generate"
)

func main() {
	// Example: Generate a random SHA-256 hash
	randomSHA256 := generate.RandomSHA256()
	fmt.Printf("Random SHA-256 Hash: %s\n", randomSHA256)

	// Example: Generate a random SHA-512 hash
	randomSHA512 := generate.RandomSHA512()
	fmt.Printf("Random SHA-512 Hash: %s\n", randomSHA512)
}
