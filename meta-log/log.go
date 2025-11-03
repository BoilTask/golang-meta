package metalog

import (
	"fmt"
	"io"
	"log/slog"
	metaconfig "meta/meta-config"
	metaerror "meta/meta-error"
	"meta/meta-flag"
	metapanic "meta/meta-panic"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

var config *Config
var logger *slog.Logger
var logFile *os.File

// Init 初始化日志系统
func Init() error {
	config = &Config{
		Stdout: true,
		Path:   "logs",
		// 每日刷新： 0 0 * * *
		// 每小时刷新： 0 0/1 ? * ?
		Cron:       "0 0/1 ? * ?",
		TimeFormat: "2006-01-02-15-04-05",
		AddSource:  false,
	}
	err := config.Parse()
	if err != nil {
		return err
	}

	err = initFileLogger()
	if err != nil {
		return err
	}

	slog.Info(
		"Log config",
		"file", metaflag.GetLogConfig(),
		"config", config,
	)

	return nil
}

func IsLogDebug() bool {
	return metaflag.IsDebug()
}

func initFileLogger() error {

	c := cron.New()
	_, err := c.AddFunc(config.Cron, createNewLogger)
	if err != nil {
		slog.Error("Error adding function to cron", "err", err)
		return err
	}
	c.Start()

	createLogger()

	return nil
}

func createNewLogger() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			metapanic.ProcessError(metaerror.Wrap(err))
			return
		}
	}
	createLogger()
}

func createLogger() {
	logDir := config.Path
	var logFile *os.File
	if logDir != "" {
		logFile = createLogFile(logDir)
	}
	var ioWriter io.Writer
	if logFile != nil && config.Stdout {
		ioWriter = io.MultiWriter(logFile, os.Stdout)
	} else if logFile != nil {
		ioWriter = logFile
	} else if config.Stdout {
		ioWriter = os.Stdout
	} else {
		ioWriter = io.Discard
	}
	options := &slog.HandlerOptions{}
	options.AddSource = config.AddSource
	if IsLogDebug() {
		options.Level = slog.LevelDebug
	} else {
		options.Level = slog.LevelInfo
	}
	if metaflag.IsDebug() {
		logger = slog.New(slog.NewTextHandler(ioWriter, options))
	} else {
		logger = slog.New(slog.NewJSONHandler(ioWriter, options))
	}
	slog.SetDefault(logger)
}

func createLogFile(logDir string) *os.File {
	currentTime := time.Now().Format(config.TimeFormat)
	fileName := fmt.Sprintf("%s_%s.log", metaconfig.GetModuleName(), currentTime)
	filePath := filepath.Join(logDir, fileName)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create log directory: %v", err))
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}
	return file
}

func GetLogger() *slog.Logger {
	return logger
}
