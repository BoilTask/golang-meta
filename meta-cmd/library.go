package metacmd

import (
	"log/slog"
	"os/exec"
	"strings"
)

func Command(name string, arg ...string) *MetaCmd {
	slog.Info("exec command", "cmd", name+" "+strings.Join(arg, " "))
	metaCmd := &MetaCmd{*exec.Command(name, arg...)}
	// 默认输出一下输入的日志，防止要求输入而卡住
	logReader := &LogReader{
		Name: name,
		Arg:  arg,
	}
	metaCmd.Stdin = logReader
	return metaCmd
}
