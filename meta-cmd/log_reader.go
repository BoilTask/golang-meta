package metacmd

import "log/slog"

type LogReader struct {
	Name string
	Arg  []string
}

func (r *LogReader) Read(p []byte) (n int, err error) {
	slog.Info("LogReader be called", "name", r.Name, "arg", r.Arg)
	return len(p), nil
}
