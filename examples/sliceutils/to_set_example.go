package main

import (
	"fmt"
	"github.com/intezya/pkglib/sliceutils"
)

func main() {
	// Example ToSetUntyped usage for slice of interface{}
	slice := []interface{}{1, "apple", 2, "orange", 1, "banana", "apple"}
	set := sliceutils.ToSetUntyped(slice)
	fmt.Println("Set of interface{} values:", set)

	// Example ToSet usage for slice of concrete type (int)
	intSlice := []int{1, 2, 2, 3, 4, 3, 5}
	intSet := sliceutils.ToSet(intSlice)
	fmt.Println("Set of int values:", intSet)

	// Example ToSet usage for slice of string
	stringSlice := []string{"apple", "banana", "apple", "orange"}
	stringSet := sliceutils.ToSet(stringSlice)
	fmt.Println("Set of string values:", stringSet)
}
