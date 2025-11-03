package engine

import (
	"flag"
	"log/slog"
	metaconfig "meta/meta-config"
	"meta/meta-flag"
	"meta/meta-log"
	metatime "meta/meta-time"
	"meta/subsystem"
	"reflect"
)

var (
	subsystems       []subsystem.Interface
	subsystemTypeMap = make(map[reflect.Type]subsystem.Interface) // 类型映射表
	subsystemNameMap = make(map[string]subsystem.Interface)       // 名称映射表
	stopChan         chan bool
)

func Init(flagFunc func(), initFuncPre func() error, initFuncEnd func() error) error {
	metatime.Init()
	metaflag.Init()
	if flagFunc != nil {
		flagFunc()
	}
	flag.Parse()
	var err error
	err = metaconfig.Init()
	if err != nil {
		return err
	}
	err = metalog.Init()
	if err != nil {
		return err
	}

	slog.Info(
		"meta config",
		"file", metaflag.GetMetaConfigFile(),
		"config", metaconfig.GetMetaConfig(),
	)

	slog.Info("engine init begin")

	if initFuncPre != nil {
		err = initFuncPre()
		if err != nil {
			return err
		}
	}

	for _, s := range subsystems {
		slog.Info("engine subsystem init begin", "subsystem", s.GetName())
		err := s.Init()
		if err != nil {
			slog.Error("engine init error", "subsystem", s.GetName(), "err", err)
			return err
		}
		slog.Info("engine subsystem init end", "subsystem", s.GetName())
	}

	if initFuncEnd != nil {
		err = initFuncEnd()
		if err != nil {
			return err
		}
	}

	slog.Info("engine init end")

	return nil
}

// Start 启动引擎
// startFuncPre 在启动引擎之前执行
// startFuncEnd 在启动引擎之后执行
// stopImmediately 是否立即结束
func Start(startFuncPre func() error, startFuncEnd func() error, stopImmediately bool) error {
	slog.Info("engine start begin")

	var err error

	if startFuncPre != nil {
		err = startFuncPre()
		if err != nil {
			return err
		}
	}

	for _, s := range subsystems {
		slog.Info("engine subsystem Start begin", "subsystem", s.GetName())
		err := s.Start()
		if err != nil {
			slog.Error("engine start error", "subsystem", s.GetName(), "err", err)
			return err
		}
		slog.Info("engine subsystem Start end", "subsystem", s.GetName())
	}

	if startFuncEnd != nil {
		err = startFuncEnd()
		if err != nil {
			return err
		}
	}

	slog.Info("engine start end")

	if !stopImmediately {
		pendingStop()
	}

	for _, s := range subsystems {
		err := s.Stop()
		if err != nil {
			slog.Error("engine stop error", "subsystem", s.GetName(), "err", err)
		}
	}

	slog.Info("engine stop")

	return nil
}

func pendingStop() {
	// 使用 channel 控制运行状态
	stopChan = make(chan bool)

	<-stopChan
}

func Stop() {
	slog.Info("engine stop manual begin")

	close(stopChan)

	slog.Info("engine stop manual end")
}

// RegisterSubsystem 应仅在initFuncPre中调用
func RegisterSubsystem(creator func() subsystem.Interface) subsystem.Interface {
	newSubsystem := creator()
	subsystems = append(subsystems, newSubsystem)
	subsystemTypeMap[reflect.TypeOf(newSubsystem)] = newSubsystem
	subsystemNameMap[newSubsystem.GetName()] = newSubsystem
	return newSubsystem
}

func GetSubsystem[T subsystem.Interface]() subsystem.Interface {
	subsystemType := reflect.TypeOf((*T)(nil)).Elem()
	thisSubsystem, found := GetSubsystemByType(subsystemType)
	if !found {
		return nil
	}
	return thisSubsystem
}

// GetSubsystemByName 根据名称返回指定的 Subsystem
func GetSubsystemByName(name string) (subsystem.Interface, bool) {
	s, found := subsystemNameMap[name]
	return s, found
}

// GetSubsystemByType 根据类型返回指定的 Subsystem
func GetSubsystemByType(subsystemType reflect.Type) (subsystem.Interface, bool) {
	s, found := subsystemTypeMap[subsystemType]
	return s, found
}
