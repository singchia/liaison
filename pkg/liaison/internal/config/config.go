package config

import "github.com/singchia/liaison/pkg/config"

// daemon related
type RLimit struct {
	Enable  bool `yaml:"enable" json:"enable"`
	NumFile int  `yaml:"nofile" json:"nofile"`
}

type PProf struct {
	Enable         bool   `yaml:"enable" json:"enable"`
	Addr           string `yaml:"addr" json:"addr"`
	CPUProfileRate int    `yaml:"cpu_profile_rate" json:"cpu_profile_rate"`
}

type Daemon struct {
	RLimit RLimit `yaml:"rlimit,omitempty" json:"rlimit"`
	PProf  PProf  `yaml:"pprof,omitempty" json:"pprof"`
}

type Manager struct {
	Listen config.Listen `yaml:"listen,omitempty" json:"listen"`
	DB     string        `yaml:"db,omitempty" json:"db"`
}

type Configuration struct {
	Daemon  Daemon  `yaml:"daemon,omitempty" json:"daemon"`
	Manager Manager `yaml:"manager,omitempty" json:"manager"`
}
