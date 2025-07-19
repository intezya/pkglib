package itertools

func Map[T, R any](values []T, f func(T) R) []R {
	res := make([]R, len(values))

	for idx, value := range values {
		res[idx] = f(value)
	}

	return res
}
