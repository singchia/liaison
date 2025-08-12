package config

import "github.com/singchia/liaison/pkg/config"

type Manager struct {
	Dial config.Dial `yaml:"dial,omitempty" json:"dial"`
}

type Configuration struct {
	Manager Manager `yaml:"manager,omitempty" json:"manager"`
}
