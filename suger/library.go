package suger

func Select[T any](q bool, a T, b T) T {
	if q {
		return a
	}
	return b
}
