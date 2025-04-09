package itertools

// Map applies a function (fn) to each element in the input slice (values)
// and returns a new slice containing the results of those function calls.
// V is the type of elements in the input slice, and U is the type of elements in the result slice.
func Map[V, U any](fn func(V) U, values []V) []U {
	result := make([]U, len(values)) // Create a result slice with the same length as the input slice
	for i, value := range values {
		result[i] = fn(value) // Apply the function to each element of the input slice
	}
	return result // Return the result slice
}

// Filter filters the elements of the input slice (values) based on a function (fn).
// The function returns only the elements for which fn returns true.
// V is the type of elements in the input and result slices.
func Filter[V any](fn func(V) bool, values []V) []V {
	result := make([]V, 0) // Initialize an empty slice to hold the filtered elements
	for _, value := range values {
		if fn(value) { // If the function returns true, append the value to the result slice
			result = append(result, value)
		}
	}
	return result // Return the filtered result slice
}
