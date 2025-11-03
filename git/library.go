package git

import (
	"log/slog"
	"meta/meta-cmd"
	"os"
	"os/exec"
	"strings"
)

func Command(name string, arg ...string) *Cmd {
	slog.Info("git command", "cmd", name+" "+strings.Join(arg, " "))
	appendArg := []string{}
	svnCmd := &Cmd{*exec.Command(name, append(appendArg, arg...)...)}
	logReader := &metacmd.LogReader{
		Name: name,
		Arg:  arg,
	}
	svnCmd.Stdin = logReader
	svnCmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	return svnCmd
}
