package metaflag

import (
	"flag"
	"meta/suger"
)

var debugFlag bool
var debugEventFlag bool
var moduleName string
var nodeName string
var metaConfigFile string
var configFile string
var logConfig string

func Init() {
	// 定义命令行参数 --debug，默认为 false
	flag.BoolVar(&debugFlag, "debug", false, "enable debug mode")
	flag.BoolVar(&debugEventFlag, "debug-event", false, "enable debug event")
	flag.StringVar(&moduleName, "module", "", "module name")
	flag.StringVar(&nodeName, "name", "", "node name")
	flag.StringVar(&metaConfigFile, "meta-config", "resource/config/meta.yaml", "meta config file")
	flag.StringVar(&configFile, "config", "resource/config/config.yaml", "config file")
	flag.StringVar(&logConfig, "log-config", "../meta/resource/config/log.yaml", "log config file")
}

func IsDebug() bool {
	return debugFlag
}

func IsDebugEvent() bool {
	return debugEventFlag
}

func HasModuleName() bool {
	return moduleName != ""
}

func GetModuleName() string {
	return suger.Select(moduleName == "", "unknown", moduleName)
}

func HasNodeName() bool {
	return nodeName != ""
}

func GetNodeName() string {
	return suger.Select(nodeName == "", "unknown", nodeName)
}

func GetMetaConfigFile() string {
	return metaConfigFile
}

func GetConfigFile() string {
	return configFile
}

func GetLogConfig() string {
	return logConfig
}
