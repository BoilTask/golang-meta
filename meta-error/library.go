package metaerror

import (
	"errors"
	"fmt"
	pkgerrors "github.com/pkg/errors"
	metaerrorcode "meta/error-code"
	metaformat "meta/meta-format"
	"reflect"
)

func GetErrorType(err error) reflect.Type {
	return reflect.TypeOf(pkgerrors.Cause(err))
}

func New(format string, param ...any) error {
	return pkgerrors.New(fmt.Sprintf(format, param...))
}

func IsErrorWithStack(err error) bool {
	type causer interface {
		Cause() error
	}
	for err != nil {
		if stack, ok := err.(interface{ StackTrace() pkgerrors.StackTrace }); ok {
			if stack.StackTrace() != nil {
				return true
			}
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return false
}

func Wrap(err error, format ...any) error {
	if err == nil {
		return nil
	}
	msg := metaformat.Format(format...)
	if IsErrorWithStack(err) {
		if msg != "" {
			return pkgerrors.WithMessage(err, msg)
		}
		return err
	}
	if msg == "" {
		msg = "meta error"
	}
	return pkgerrors.Wrap(err, msg)
}

func NewCode[T metaerrorcode.Numeric](code T, format ...any) error {
	return &MetaError{
		Code: int(code),
		Err:  pkgerrors.New(metaformat.Format(format...)),
	}
}

func WrapCode[T metaerrorcode.Numeric](err error, code T, format ...any) error {
	if err == nil {
		return nil
	}
	msg := metaformat.Format(format...)
	var metaErr *MetaError
	if errors.As(err, &metaErr) {
		metaErr.Code = int(code)
		return pkgerrors.WithMessage(metaErr, msg)
	}
	if _, ok := err.(interface{ StackTrace() pkgerrors.StackTrace }); ok {
		return &MetaError{
			Code: int(code),
			Err:  pkgerrors.WithMessage(err, msg),
		}
	}
	return &MetaError{
		Code: int(code),
		Err:  pkgerrors.Wrap(err, msg),
	}
}

func WrapFeishu(err interface{}, format ...any) error {
	if err == nil {
		return nil
	}
	type ErrorInfo interface {
		ErrorResp() string
	}
	switch v := err.(type) {
	case ErrorInfo:
		return Wrap(New(v.ErrorResp()), format...)
	case error:
		return Wrap(err.(error), format...)
	default:
		return New("feishu err:%s message:%s", fmt.Sprint(err), metaformat.Format(format...))
	}
}

func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	if n == 1 {
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, Wrap(err))
		}
	}
	return e
}

func GetErrorCodeFromError(err error) int {
	code := int(metaerrorcode.UnknownError)
	var metaError *MetaError
	ok := errors.As(err, &metaError)
	if ok {
		code = metaError.Code
	}
	return code
}
