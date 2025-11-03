package metalog

import (
	metaerror "meta/meta-error"
	"meta/meta-flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Stdout     bool   `yaml:"stdout"`
	Path       string `yaml:"path"`
	Cron       string `yaml:"cron"`
	TimeFormat string `yaml:"time-format"`
	AddSource  bool   `yaml:"add-source"`
}

func (config *Config) Parse() error {
	configFile := metaflag.GetLogConfig()
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return metaerror.Wrap(err, "failed to read log config file: %s", configFile)
	}
	return yaml.Unmarshal(yamlFile, config)
}
