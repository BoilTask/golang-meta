package metaerror

import (
	"errors"
	"fmt"
	pkgerrors "github.com/pkg/errors"
	"io"
)

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	var message string
	for _, err := range e.errs {
		message += err.Error() + "\n"
	}
	return message
}

func (e *joinError) Format(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, fmt.Sprintf("[%d]", len(e.errs)))
	for _, err := range e.errs {
		if formatter, ok := err.(fmt.Formatter); ok {
			// 如果 f.errs 支持格式化，递归调用其 Format 方法
			formatter.Format(s, verb)
		} else {
			// 否则，使用默认的格式化规则
			_, _ = io.WriteString(s, fmt.Sprintf("%"+string(verb), e.errs))
		}
	}
}

func (e *joinError) StackTrace() pkgerrors.StackTrace {
	stacks := make(pkgerrors.StackTrace, 0)
	for _, err := range e.errs {
		if err != nil {
			stack, ok := err.(interface{ StackTrace() pkgerrors.StackTrace })
			if ok {
				stacks = append(stacks, stack.StackTrace()...)
			}
		}
	}
	return stacks
}

func (e *joinError) Cause() error {
	var fallbackError error
	for _, err := range e.errs {
		if err != nil {
			cause := pkgerrors.Cause(err)
			if !errors.Is(err, cause) {
				return cause
			}
			if fallbackError == nil {
				fallbackError = err
			}
		}
	}
	return fallbackError
}
