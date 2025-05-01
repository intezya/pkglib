package itertools

func GetMapKeys[K, V comparable](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}

func GetMapValues[K, V comparable](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, value := range m {
		result = append(result, value)
	}
	return result
}

func MapMapKeys[K, V, U comparable](m map[K]V, fn func(K) U) []U {
	keys := GetMapKeys(m)
	return Map(keys, fn)
}

func MapMapValues[K, V, U comparable](m map[K]V, fn func(V) U) []U {
	values := GetMapValues(m)
	return Map(values, fn)
}
