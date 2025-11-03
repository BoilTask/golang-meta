package metapanic

import metaerror "meta/meta-error"

var ProcessPanicCallback func(name string, err error, format ...any)
var ProcessErrorCallback func(err error, format ...any)

func ProcessError(err error, format ...any) {
	if err == nil {
		return
	}
	if ProcessErrorCallback != nil {
		ProcessErrorCallback(err, format...)
	} else {
		LogError(err, format...)
	}
}

func ProcessPanic(name string, err interface{}, format ...any) {
	if err == nil {
		return
	}
	var realErr error
	switch v := err.(type) {
	case error:
		realErr = v
	case string:
		realErr = metaerror.New(v)
	default:
		realErr = metaerror.New("unknown panic: %+v", v)
	}
	if ProcessPanicCallback != nil {
		ProcessPanicCallback(name, realErr, format...)
	} else {
		LogPanic(name, realErr, format...)
	}
}
