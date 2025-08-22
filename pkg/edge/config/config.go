package config

import "github.com/singchia/liaison/pkg/config"

type Auth struct {
	AccessKey string `yaml:"access_key,omitempty" json:"access_key"`
	SecretKey string `yaml:"secret_key,omitempty" json:"secret_key"`
}

type Manager struct {
	Dial config.Dial `yaml:"dial,omitempty" json:"dial"`
	Auth Auth        `yaml:"auth,omitempty" json:"auth"`
}

type Configuration struct {
	Manager Manager `yaml:"manager,omitempty" json:"manager"`
}
