package main

import (
	"fmt"
	"github.com/intezya/pkglib/configloader"
	"log"
)

// go run ... --env-file=.env
func main() {
	configloader.LoadEnv()

	// Example for getting an environment variable
	apiKey, err := configloader.GetEnv("API_KEY", "default-api-key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("API Key:", apiKey)

	// Example for getting an integer environment variable
	port, err := configloader.GetEnvInt("PORT", 8080)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port:", port)
}
