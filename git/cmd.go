package git

import (
	"log/slog"
	"os/exec"
	"slices"
)

type Cmd struct {
	exec.Cmd
}

var excludeArg = []string{"--password"}

func (c *Cmd) Run() error {
	var logArg []string
	for i := 0; i < len(c.Args); i++ {
		arg := c.Args[i]
		// 如果包括的跳过参数，那么跳过这个参数和他后面的值
		if slices.Contains(excludeArg, arg) {
			i += 1
			continue
		}
		logArg = append(logArg, arg)
	}
	slog.Info("command run", "dir", c.Dir, "path", c.Path, "arg", logArg)
	return c.Cmd.Run()
}
