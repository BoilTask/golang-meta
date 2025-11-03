package metafeishu

import (
	"context"
	"fmt"
	"log/slog"
	metalog "meta/meta-log"
)

type Logger struct {
}

func IsSwitchDebugToInfo() bool {
	// 如果开启着日志的Debug输出，没必要转换
	if metalog.IsLogDebug() {
		return false
	}
	// 目前访问量较小，因此将飞书的Debug日志视为Info日志，方便监听
	return true
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	if IsSwitchDebugToInfo() {
		slog.InfoContext(ctx, "Feishu Debug", "args", fmt.Sprint(args...))
	} else {
		slog.DebugContext(ctx, "Feishu Debug", "args", fmt.Sprint(args...))
	}
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	slog.InfoContext(ctx, "Feishu Info", "args", fmt.Sprint(args...))
}

func (l *Logger) Warn(ctx context.Context, args ...interface{}) {
	slog.WarnContext(ctx, "Feishu Warn", "args", fmt.Sprint(args...))
}

func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	slog.ErrorContext(ctx, "Feishu Error", "args", fmt.Sprint(args...))
}

func NewLogger() *Logger {
	return &Logger{}
}
