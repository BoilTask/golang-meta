package svn

import (
	"log/slog"
	"meta/meta-cmd"
	"os/exec"
	"strings"
)

func Command(username string, password string, name string, arg ...string) *Cmd {
	slog.Info("svn command", "cmd", name+" "+strings.Join(arg, " "))
	appendArg := []string{"--username", username, "--password", password, "--non-interactive"}
	svnCmd := &Cmd{*exec.Command(name, append(appendArg, arg...)...)}
	logReader := &metacmd.LogReader{
		Name: name,
		Arg:  arg,
	}
	svnCmd.Stdin = logReader
	return svnCmd
}

func CommandWithoutAuth(name string, arg ...string) *Cmd {
	slog.Info("svn command", "cmd", name+" "+strings.Join(arg, " "))
	appendArg := []string{"--non-interactive"}
	svnCmd := &Cmd{*exec.Command(name, append(appendArg, arg...)...)}
	logReader := &metacmd.LogReader{
		Name: name,
		Arg:  arg,
	}
	svnCmd.Stdin = logReader
	return svnCmd
}
