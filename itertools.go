package pkglib

type Pair[T, U any] struct {
	First  T
	Second U
}

func Repeat[V any](value V, count int) []V {
	result := make([]V, count)
	for i := 0; i < count; i++ {
		result[i] = value
	}
	return result
}

func Map[V, U any](fn func(V) U, values []V) []U {
	result := make([]U, len(values))
	for i, value := range values {
		result[i] = fn(value)
	}
	return result
}

func Filter[V any](fn func(V) bool, values []V) []V {
	result := make([]V, 0)
	for _, value := range values {
		if fn(value) {
			result = append(result, value)
		}
	}
	return result
}

func Reduce[V any](fn func(V, V) V, values []V) V {
	if len(values) == 0 {
		panic("Reduce on empty slice")
	}
	result := values[0]
	for _, value := range values[1:] {
		result = fn(result, value)
	}
	return result
}

func Zip[T, U any](values1 []T, values2 []U) []Pair[T, U] {
	minLen := len(values1)
	if len(values2) < minLen {
		minLen = len(values2)
	}

	result := make([]Pair[T, U], minLen)
	for i := 0; i < minLen; i++ {
		result[i] = Pair[T, U]{values1[i], values2[i]}
	}
	return result
}
