package itertools

import (
	"iter"
	"slices"
)

type Iterator[T any] struct {
	impl iter.Seq[T]
}

func From[T any](data []T) *Iterator[T] {
	return &Iterator[T]{
		impl: func(yield func(T) bool) {
			for _, value := range data {
				if !yield(value) {
					return
				}
			}
		},
	}
}

func (i *Iterator[T]) Collect() []T {
	var data []T

	for value := range i.impl {
		data = append(data, value)
	}

	return data
}

func (i *Iterator[T]) Each(fn func(T)) {
	for value := range i.impl {
		fn(value)
	}
}

func (i *Iterator[T]) Filter(predicate func(T) bool) *Iterator[T] {
	copied := i.impl

	i.impl = func(yield func(T) bool) {
		for value := range copied {
			if predicate(value) {
				if !yield(value) {
					return
				}
			}
		}
	}

	return i
}

func (i *Iterator[T]) Map(fn func(T) T) *Iterator[T] {
	copied := i.impl

	i.impl = func(yield func(T) bool) {
		for value := range copied {
			value = fn(value)
			if !yield(value) {
				return
			}
		}
	}

	return i
}

func (i *Iterator[T]) Reverse() *Iterator[T] {
	data := i.Collect()

	slices.Reverse(data)

	return From(data)
}

// Count consumes the iterator and returns the number of elements.
func (i *Iterator[T]) Count() int {
	c := 0

	for range i.impl {
		c++
	}

	return c
}

// CountWithPredicate consumes the iterator and returns the number of elements that satisfies predicate.
func (i *Iterator[T]) CountWithPredicate(predicate func(T) bool) int {
	return i.Filter(predicate).Count()
}

// Example
//func main() {
//	From([]int{1, 2, 3, 4}).
//		Reverse().
//		Map(func(x int) int { return x * x }).
//		Filter(func(x int) bool { return x%2 == 0 }).
//		Each(func(x int) { fmt.Println(x) })
//}
