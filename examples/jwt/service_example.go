package main

import (
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/intezya/pkglib/jwt"
	"log"
	"time"
)

func main() {
	// Example configuration for the JWT service
	config := jwt.NewConfiguration([]byte("mySecretKey"), "myApp", time.Hour)

	// Create a new JWT service
	jwtService := jwt.New(config)

	// Generate a token with authentication data
	authData := jwt.AuthenticationData{"userId": "1234", "role": "admin"}
	token, err := jwtService.Generate(authData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generated Token:", token)

	// Validate the token and extract authentication data
	decodedData, err := jwtService.Validate(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decoded Authentication Data:", decodedData)

	// Validate the token with custom options
	decodedData, err = jwtService.Validate(token, jwt2.WithIssuedAt(), jwt2.WithStrictDecoding())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decoded Authentication Data:", decodedData)
}
