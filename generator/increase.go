package generator

import (
	"golang.org/x/exp/constraints"
	"sync"
)

type IncreaseGenerator[T constraints.Integer | constraints.Float] struct {
	mutex   sync.RWMutex
	counter T
	step    T
}

func NewIncreaseGenerator[T constraints.Integer | constraints.Float](first T, step T) *IncreaseGenerator[T] {
	return &IncreaseGenerator[T]{
		counter: first - step,
		step:    step,
	}
}

func (gen *IncreaseGenerator[T]) Next() T {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()
	gen.counter += gen.step
	return gen.counter
}
