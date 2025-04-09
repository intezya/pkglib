package main

import (
	"fmt"
	"github.com/intezya/pkglib/itertools"
)

func main() {
	// Example: Map - Convert each number in the slice to its square
	numbers := []int{1, 2, 3, 4, 5}
	squares := itertools.Map(
		func(n int) int {
			return n * n
		}, numbers,
	)
	fmt.Println("Squares:", squares) // Output: Squares: [1 4 9 16 25]

	// Example: Filter - Keep only even numbers from the slice
	evens := itertools.Filter(
		func(n int) bool {
			return n%2 == 0
		}, numbers,
	)
	fmt.Println("Evens:", evens) // Output: Evens: [2 4]

}
