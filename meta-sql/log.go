package metasql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	gormlogger.Config
}

// LogMode log mode
func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		slog.InfoContext(ctx, "gorm", "msg", fmt.Sprintf(msg, data...), "stack", utils.FileWithLineNum())
	}
}

// Warn print warn messages
func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		slog.WarnContext(ctx, "gorm", "msg", fmt.Sprintf(msg, data...), "stack", utils.FileWithLineNum())
	}
}

// Error print error messages
func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		slog.ErrorContext(ctx, "gorm", "msg", fmt.Sprintf(msg, data...), "stack", utils.FileWithLineNum())
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(
		err,
		gormlogger.ErrRecordNotFound,
	) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			slog.ErrorContext(
				ctx, "gorm",
				"err",
				err,
				"elapsed",
				float64(elapsed.Nanoseconds())/1e6,
				"rows",
				"-",
				"sql",
				sql,
				"stack",
				utils.FileWithLineNum(),
			)
		} else {
			slog.ErrorContext(
				ctx, "gorm",
				"err",
				err,
				"elapsed",
				float64(elapsed.Nanoseconds())/1e6,
				"rows",
				rows,
				"sql",
				sql,
				"stack",
				utils.FileWithLineNum(),
			)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			slog.WarnContext(
				ctx, "gorm",
				"slowLog",
				slowLog,
				"elapsed",
				float64(elapsed.Nanoseconds())/1e6,
				"rows",
				"-",
				"sql",
				sql,
				"stack",
				utils.FileWithLineNum(),
			)
		} else {
			slog.WarnContext(
				ctx, "gorm",
				"slowLog",
				slowLog,
				"elapsed",
				float64(elapsed.Nanoseconds())/1e6,
				"rows",
				rows,
				"sql",
				sql,
				"stack",
				utils.FileWithLineNum(),
			)
		}
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			slog.InfoContext(
				ctx, "gorm",
				"elapsed",
				float64(elapsed.Nanoseconds())/1e6,
				"rows",
				"-",
				"sql",
				sql,
				"stack",
				utils.FileWithLineNum(),
			)
		} else {
			slog.InfoContext(
				ctx, "gorm",
				"elapsed",
				float64(elapsed.Nanoseconds())/1e6,
				"rows",
				rows,
				"sql",
				sql,
				"stack",
				utils.FileWithLineNum(),
			)
		}
	}
}
