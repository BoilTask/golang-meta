package routine

import (
	metapanic "meta/meta-panic"
	"time"
)

var SafeGoCallback func(name string, err interface{})

func Protect(name string) {
	if err := recover(); err != nil {
		if SafeGoCallback != nil {
			SafeGoCallback(name, err)
		} else {
			metapanic.ProcessPanic(name, err)
		}
	}
}

func ProtectWithCallback(name string, callback func(name string, err interface{})) {
	if err := recover(); err != nil {
		if callback != nil {
			callback(name, err)
		} else {
			if SafeGoCallback != nil {
				SafeGoCallback(name, err)
			} else {
				metapanic.ProcessPanic(name, err)
			}
		}
	}
}

func SafeGo(name string, task func() error) {
	go func() {
		func() {
			defer Protect(name)
			err := task()
			if err != nil {
				metapanic.ProcessError(err, name)
				return
			}
		}()
	}()
}

func SafeGoWithRestart(name string, task func() error) {
	go func() {
		needBreak := false
		for {
			func() {
				defer Protect(name)
				err := task()
				needBreak = true
				if err != nil {
					metapanic.ProcessError(err, name)
					return
				}
			}()
			if needBreak {
				break
			}
			time.Sleep(time.Second * 1) // 避免 goroutine 频繁崩溃导致 CPU 过载
		}
	}()
}

func SafeGoWithCallback(name string, task func() error, callback func(name string, err interface{})) {
	go func() {
		defer ProtectWithCallback(name, callback)
		err := task()
		if err != nil {
			metapanic.ProcessError(err, name)
			return
		}
	}()
}
