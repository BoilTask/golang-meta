package metaerrorcode

type Numeric interface {
	~int
}

type ErrorCode int

const (
	Success         ErrorCode = 0
	CommonError     ErrorCode = 1000
	PanicError      ErrorCode = 1001
	UnknownError    ErrorCode = 1002
	TooManyRequests ErrorCode = 1003 // 太过频繁
)
