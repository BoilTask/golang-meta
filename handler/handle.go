package handler

import "context"

type Handler[T any] interface {
	// Init 初始化
	Init() error
	// IsShouldProcess 判断是否应该处理
	IsShouldProcess(ctx context.Context, e T) bool
	// DoProcess 处理事件，返回是否继续处理
	DoProcess(ctx context.Context, e T) (bool, error)
}
