package queue

import (
	"container/heap"
)

// Item 定义优先队列的元素类型
type Item[T any] struct {
	Value T   // 存储的值
	Index int // 用于内部管理的索引
}

// PriorityQueue 定义优先队列
type PriorityQueue[T any] struct {
	items    []*Item[T]        // 存储队列中的元素
	lessFunc func(a, b T) bool // 比较函数，用于决定优先级
}

// NewPriorityQueue 创建一个优先队列
func NewPriorityQueue[T any](lessFunc func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items:    []*Item[T]{},
		lessFunc: lessFunc,
	}
}

// Len 返回队列长度
func (pq *PriorityQueue[T]) Len() int {
	return len(pq.items)
}

// Less 定义优先级比较逻辑
func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return pq.lessFunc(pq.items[i].Value, pq.items[j].Value)
}

// Swap 交换两个元素的位置
func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].Index = i
	pq.items[j].Index = j
}

// Push 将新元素加入队列
func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*Item[T])
	item.Index = len(pq.items)
	pq.items = append(pq.items, item)
}

// Pop 移除优先级最高的元素
func (pq *PriorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.Index = -1 // 清除索引
	pq.items = old[0 : n-1]
	return item
}

// Add 添加新元素
func (pq *PriorityQueue[T]) Add(value T) {
	heap.Push(pq, &Item[T]{Value: value})
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.Len() == 0
}

// PopHighestPriority 移除并返回优先级最高的元素
func (pq *PriorityQueue[T]) PopHighestPriority() T {
	item := heap.Pop(pq).(*Item[T])
	return item.Value
}

func (pq *PriorityQueue[T]) TopHighestPriority() T {
	return pq.items[0].Value
}
