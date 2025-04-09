package main

import (
	"fmt"
	"github.com/intezya/pkglib/crypto"
	"log"
)

func main() {
	// Example: Generate a salt of 16 bytes
	salt, err := crypto.Salt(16)
	if err != nil {
		log.Fatalf("Error generating salt: %v", err)
	}
	fmt.Printf("Generated Salt: %x\n", salt)
}
