package main

import (
	"fmt"
	"github.com/intezya/pkglib/generate"
)

func main() {
	// Example: Generate a random string of length 16 using the AllCharset (letters, digits, symbols)
	randomStr := generate.RandomString(16, generate.Charset.AllCharset)
	fmt.Printf("Random String (All Charset): %s\n", randomStr)
}
