package metaerror

import (
	"fmt"
	pkgerrors "github.com/pkg/errors"
	"io"
)

type MetaError struct {
	Code int
	Err  error
}

func (e *MetaError) Unwrap() error {
	return e.Err
}

func (e *MetaError) Error() string {
	return e.Err.Error()
}

func (e *MetaError) Format(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, fmt.Sprintf("[%d]", e.Code))
	if formatter, ok := e.Err.(fmt.Formatter); ok {
		// 如果 f.errs 支持格式化，递归调用其 Format 方法
		formatter.Format(s, verb)
	} else {
		// 否则，使用默认的格式化规则
		_, _ = io.WriteString(s, fmt.Sprintf("%"+string(verb), e.Err))
	}
}

func (e *MetaError) StackTrace() pkgerrors.StackTrace {
	stack, ok := e.Err.(interface{ StackTrace() pkgerrors.StackTrace })
	if !ok {
		return nil
	}
	return stack.StackTrace()
}

func (e *MetaError) Cause() error {
	return pkgerrors.Cause(e.Err)
}
