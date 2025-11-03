package singleton

import "sync"

type Singleton[T any] struct {
	once     sync.Once
	instance *T
}

// GetInstance 返回单例实例
func (s *Singleton[T]) GetInstance(creator func() *T) *T {
	s.once.Do(func() {
		inst := creator()
		s.instance = inst
	})
	return s.instance
}
