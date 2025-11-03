package metacmd

import (
	"log/slog"
	"os/exec"
	"slices"
)

type MetaCmd struct {
	exec.Cmd
}

func (c *MetaCmd) Run() error {
	slog.Info("command run", "dir", c.Dir, "path", c.Path, "arg", c.Args)
	return c.Cmd.Run()
}

func (c *MetaCmd) RunExcludeLog(excludeArg ...string) error {
	var logArg []string
	for _, arg := range c.Args {
		if !slices.Contains(excludeArg, arg) {
			logArg = append(logArg, arg)
		}
	}
	slog.Info("command run", "dir", c.Dir, "path", c.Path, "arg", logArg)
	return c.Cmd.Run()
}
