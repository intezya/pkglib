package pkglib

type itertools struct{}

func (itertools) Repeat(value interface{}, count int) []interface{} {
	result := make([]interface{}, count)
	for i := 0; i < count; i++ {
		result[i] = value
	}
	return result
}

func (itertools) Map(fn func(interface{}) interface{}, values []interface{}) []interface{} {
	result := make([]interface{}, len(values))
	for i, value := range values {
		result[i] = fn(value)
	}
	return result
}

func (itertools) Filter(fn func(interface{}) bool, values []interface{}) []interface{} {
	result := make([]interface{}, 0)
	for _, value := range values {
		if fn(value) {
			result = append(result, value)
		}
	}
	return result
}

func (itertools) Reduce(fn func(interface{}, interface{}) interface{}, values []interface{}) interface{} {
	result := values[0]
	for _, value := range values[1:] {
		result = fn(result, value)
	}
	return result
}

func (itertools) Zip(values1 []interface{}, values2 []interface{}) []interface{} {
	result := make([]interface{}, len(values1))
	for i, value1 := range values1 {
		result[i] = []interface{}{value1, values2[i]}
	}
	return result
}

var IterTools itertools
