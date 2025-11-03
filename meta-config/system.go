package metaconfig

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	metaerror "meta/meta-error"
	"meta/meta-flag"
	metaformat "meta/meta-format"
	"meta/subsystem"
	"os"
)

type Subsystem struct {
	subsystem.Subsystem
}

var configInterface ConfigInterface

func (s *Subsystem) GetName() string {
	return "Config"
}

func (s *Subsystem) InitConfig(config ConfigInterface) error {

	configFile := metaflag.GetConfigFile()

	slog.Info("Config init", "configFile", configFile)

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return metaerror.Wrap(err, "Unmarshal config file error")
	}

	configInterface = config

	if metaflag.IsDebug() {
		slog.Info(
			"Init Config", "config",
			metaformat.FormatByJson(config),
		)
	}

	return nil
}

func GetConfig() ConfigInterface {
	return configInterface
}
