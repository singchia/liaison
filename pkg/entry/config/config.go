package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
}

var Conf *Configuration

func Init() (*Configuration, error) {
	configPath := os.Getenv("ENTRY_CONFIG_PATH")
	if configPath == "" {
		configPath = "entry.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		Conf = &Configuration{}
		return Conf, nil
	}

	conf := &Configuration{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	Conf = conf
	return Conf, nil
}
