package retry

import (
	"log/slog"
	metaerror "meta/meta-error"
	"time"
)

func TryRetry(name string, maxCount int, f func(int) bool) error {
	for i := 0; i < maxCount; i++ {
		if f(i) {
			return nil
		}
		slog.Info("wait retry",
			"name", name,
			"current index", i,
			"max count", maxCount,
		)
	}
	return metaerror.New("retry failed, name:%s, max count:%d", name, maxCount)
}

func TryRetrySleep(name string, maxCount int, sleep time.Duration, f func(int) bool) error {
	for i := 0; i < maxCount; i++ {
		if f(i) {
			if i > 0 {
				slog.Info("retry sleep success",
					"name", name,
					"current index", i,
					"max count", maxCount,
					"sleep", sleep,
				)
			}
			return nil
		}
		slog.Info("wait retry sleep",
			"name", name,
			"current index", i,
			"max count", maxCount,
			"sleep", sleep,
		)
		time.Sleep(sleep)
	}
	return metaerror.New("retry sleep failed, name:%s, max count:%d, sleep:%s",
		name, maxCount, sleep,
	)
}

func TryRetryDynamicSleep(name string, maxCount int, f func(int) *time.Duration) error {
	for i := 0; i < maxCount; i++ {
		sleep := f(i)
		if sleep == nil {
			if i > 0 {
				slog.Info("retry sleep success",
					"name", name,
					"current index", i,
					"max count", maxCount,
					"sleep", sleep,
				)
			}
			return nil
		}
		slog.Info("wait retry sleep",
			"name", name,
			"current index", i,
			"max count", maxCount,
			"sleep", sleep,
		)
		time.Sleep(*sleep)
	}
	return metaerror.New("retry sleep failed, name:%s, max count:%d", name, maxCount)
}

func TryRetryWhenErr(name string, maxCount int, f func(int) error) error {
	var err error
	for i := 0; i < maxCount; i++ {
		err = f(i)
		if err == nil {
			return nil
		}
		slog.Info("wait retry when err",
			"name", name,
			"current index", i,
			"max count", maxCount,
			"err", err,
		)
	}
	return err
}

func TryRetryWhenErrSleep(name string, maxCount int, sleep time.Duration, f func(int) error) error {
	var err error
	for i := 0; i < maxCount; i++ {
		err = f(i)
		if err == nil {
			return nil
		}
		slog.Info("wait retry when err sleep",
			"name", name,
			"current index", i,
			"max count", maxCount,
			"sleep", sleep,
			"err", err,
		)
		time.Sleep(sleep)
	}
	return err
}
