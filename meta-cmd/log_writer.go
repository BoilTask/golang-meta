package metacmd

import "log/slog"

type LogWriter struct {
	Name string
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	slog.Info("LogWriter", "name", w.Name, "output", string(p))
	return len(p), nil
}
