package itertools

// Map applies a function (fn) to each element in the input slice (values)
// and returns a new slice containing the results of those function calls.
// V is the type of elements in the input slice, and U is the type of elements in the result slice.
func Map[V, U any](values []V, fn func(V) U) []U {
	result := make([]U, len(values))
	for i, value := range values {
		result[i] = fn(value)
	}
	return result
}

// Filter filters the elements of the input slice (values) based on a function (fn).
// The function returns only the elements for which fn returns true.
// V is the type of elements in the input and result slices.
func Filter[V any](fn func(V) bool, values []V) []V {
	result := make([]V, 0)
	for _, value := range values {
		if fn(value) {
			result = append(result, value)
		}
	}
	return result
}
