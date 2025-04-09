package sliceutils

// Set represents a set, where each element is stored as a key in a map with boolean values (true).
// The type parameter T must be comparable to ensure it can be used as a map key.
type Set[T comparable] map[T]bool

// ToSetUntyped converts a slice of any type (interface{}) into a set of interface{} values.
// It ensures each value in the slice is unique by using a map where the key is the value and the value is true.
// This function returns a Set[interface{}] containing only unique elements from the provided slice.
func ToSetUntyped(slice []interface{}) Set[interface{}] {
	set := make(Set[interface{}], len(slice)) // Create a new Set with an initial capacity based on slice length
	for _, value := range slice {
		set[value] = true // Store each value in the map (duplicates will be overwritten)
	}
	return set
}

// ToSet converts a slice of any comparable type (T) into a set of unique elements of type T.
// It ensures each element in the slice is unique by using a map where the key is the element and the value is true.
// This function returns a Set[T] containing only unique elements from the provided slice.
func ToSet[T comparable](slice []T) Set[T] {
	set := make(Set[T], len(slice)) // Create a new Set with an initial capacity based on slice length
	for _, value := range slice {
		set[value] = true // Store each value in the map (duplicates will be overwritten)
	}
	return set
}
