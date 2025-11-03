package metapanic

import (
	"fmt"
	"log/slog"
	"meta/host"
	metaconfig "meta/meta-config"
	metaerror "meta/meta-error"
	metaformat "meta/meta-format"
)

func LogError(err error, format ...any) string {
	message := fmt.Sprintf("error [%s] [%s] %s", metaconfig.GetModuleName(), metaconfig.GetNodeName(), host.GetLocalIp())
	content := metaformat.Format(format...)
	if content != "" {
		message = fmt.Sprintf("%s\n%s", message, content)
	}
	errStr := fmt.Sprintf("[%s] %+v", metaerror.GetErrorType(err), metaerror.Wrap(err))
	message = fmt.Sprintf("%s\n%s", message, errStr)
	slog.Error("Error", "message", message)
	return message
}

func LogPanic(name string, err error, format ...any) string {
	message := fmt.Sprintf("panic [%s] [%s] %s", metaconfig.GetModuleName(), metaconfig.GetNodeName(), host.GetLocalIp())
	message = fmt.Sprintf("%s\n%s", message, name)
	content := metaformat.Format(format...)
	if content != "" {
		message = fmt.Sprintf("%s\n%s", message, content)
	}
	errStr := fmt.Sprintf("[%s] %+v", metaerror.GetErrorType(err), metaerror.Wrap(err))
	message = fmt.Sprintf("%s\n%s", message, errStr)
	slog.Error("Panic", "message", message)
	return message
}
