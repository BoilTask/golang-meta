package metaconfig

import (
	"gopkg.in/yaml.v3"
	metaerror "meta/meta-error"
	"meta/meta-flag"
	"os"
)

var metaConfig Config

func Init() error {
	configFile := metaflag.GetMetaConfigFile()
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &metaConfig)
	if err != nil {
		return metaerror.Wrap(err, "Unmarshal config file error")
	}
	return nil
}

func GetMetaConfig() *Config {
	return &metaConfig
}

func GetModuleName() string {
	if metaflag.HasModuleName() {
		return metaflag.GetModuleName()
	}
	return GetMetaConfig().Module.Name
}

func GetNodeName() string {
	if metaflag.HasNodeName() {
		return metaflag.GetNodeName()
	}
	return GetMetaConfig().Node.Name
}
