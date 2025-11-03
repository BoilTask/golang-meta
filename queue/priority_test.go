package queue

import (
	"testing"
)

// 测试整数优先队列的行为
func TestPriorityQueueWithInt(t *testing.T) {
	// 定义整数优先级比较函数（小值优先）
	intLessFunc := func(a, b int) bool {
		return a < b
	}

	// 创建优先队列
	pq := NewPriorityQueue(intLessFunc)

	// 添加元素
	pq.Add(10)
	pq.Add(5)
	pq.Add(20)

	// 期望的顺序
	expectedOrder := []int{5, 10, 20}

	for i, expected := range expectedOrder {
		if pq.Len() == 0 {
			t.Fatalf("Queue ended prematurely at index %d", i)
		}
		actual := pq.PopHighestPriority()
		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	}
}

// 测试复杂类型优先队列
func TestPriorityQueueWithTasks(t *testing.T) {
	// 定义任务类型
	type Task struct {
		ID       int
		Priority int
		Name     string
	}

	// 定义任务优先级比较函数（按优先级升序）
	taskLessFunc := func(a, b Task) bool {
		return a.Priority < b.Priority
	}

	// 创建优先队列
	pq := NewPriorityQueue(taskLessFunc)

	// 添加任务
	pq.Add(Task{ID: 1, Priority: 10, Name: "Task A"})
	pq.Add(Task{ID: 2, Priority: 5, Name: "Task B"})
	pq.Add(Task{ID: 3, Priority: 20, Name: "Task C"})

	// 期望的顺序
	expectedOrder := []Task{
		{ID: 2, Priority: 5, Name: "Task B"},
		{ID: 1, Priority: 10, Name: "Task A"},
		{ID: 3, Priority: 20, Name: "Task C"},
	}

	for i, expected := range expectedOrder {
		if pq.Len() == 0 {
			t.Fatalf("Queue ended prematurely at index %d", i)
		}
		actual := pq.PopHighestPriority()
		if actual != expected {
			t.Errorf("Expected %+v, got %+v", expected, actual)
		}
	}
}

// 测试空队列行为
func TestEmptyQueue(t *testing.T) {
	// 定义整数优先级比较函数（小值优先）
	intLessFunc := func(a, b int) bool {
		return a < b
	}

	// 创建优先队列
	pq := NewPriorityQueue(intLessFunc)

	// 检查队列为空
	if pq.Len() != 0 {
		t.Fatalf("Expected empty queue, but got length %d", pq.Len())
	}
}
