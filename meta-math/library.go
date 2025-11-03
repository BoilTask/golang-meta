package metamath

import (
	. "golang.org/x/exp/constraints"
	"math"
)

func Abs[T Signed | Float](a T) T {
	if a < T(0) {
		return -a
	}
	return a
}

// Clamp 将值限制在[min, max]范围内
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Min[T Ordered](a T, b ...T) T {
	minVal := a
	for _, v := range b {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func Max[T Ordered](a T, b ...T) T {
	maxVal := a
	for _, v := range b {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func Log(base float64, value float64) float64 {
	return math.Log(value) / math.Log(base)
}

func LogE(value float64) float64 {
	return math.Log(value)
}
