package set

import "sort"

type Set[T comparable] map[T]struct{}

func New[T comparable](size ...int) Set[T] {
	if len(size) == 0 {
		return make(Set[T])
	}
	return make(Set[T], size[0])
}

func FromSlice[T comparable](slice []T) Set[T] {
	result := New[T]()
	for _, value := range slice {
		result.Add(value)
	}
	return result
}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	result := New[T]()
	for value := range s {
		result.Add(value)
	}
	for value := range other {
		result.Add(value)
	}
	return result
}

func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for value := range s {
		result = append(result, value)
	}
	return result
}

func (s Set[T]) ToSortSlice(less func(a, b T) bool) []T {
	result := make([]T, 0, len(s))
	for value := range s {
		result = append(result, value)
	}
	sort.Slice(
		result, func(i, j int) bool {
			return less(result[i], result[j])
		},
	)
	return result
}

func (s Set[T]) Foreach(f func(value *T) bool) {
	for value := range s {
		if !f(&value) {
			break
		}
	}
}

func (s Set[T]) ForeachCopy(f func(value T) bool) {
	for value := range s {
		if !f(value) {
			break
		}
	}
}
