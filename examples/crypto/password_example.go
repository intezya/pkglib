package main

import (
	"fmt"
	"github.com/intezya/pkglib/crypto"
	"log"
)

func main() {
	// Example: Generate a salt for Argon2
	salt, err := crypto.Salt(16)
	if err != nil {
		log.Fatalf("Error generating salt: %v", err)
	}

	// Example: Hash a password using Argon2
	password := "securepassword"
	hash := crypto.HashArgon2(password, salt, nil)
	fmt.Printf("Argon2 Hash: %s\n", hash)

	// Example: Verify the password against the hash
	isValid := crypto.VerifyArgon2(password, hash, nil)
	fmt.Printf("Password valid: %v\n", isValid)
}
