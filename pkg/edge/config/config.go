package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/singchia/liaison/pkg/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

var (
	Conf      *Configuration
	RotateLog *lumberjack.Logger

	h           bool
	showFingerprint bool
	file        string
	defaultFile string = "./liaison-edge.yaml"
)

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

type Auth struct {
	AccessKey string `yaml:"access_key,omitempty" json:"access_key"`
	SecretKey string `yaml:"secret_key,omitempty" json:"secret_key"`
}

type Manager struct {
	Dial config.Dial `yaml:"dial,omitempty" json:"dial"`
	Auth Auth        `yaml:"auth,omitempty" json:"auth"`
}

type Log struct {
	Level    string `yaml:"level"`
	File     string `yaml:"file"`
	MaxSize  int    `yaml:"maxsize"`
	MaxRolls int    `yaml:"maxrolls"`
}

type Configuration struct {
	Daemon  Daemon  `yaml:"daemon,omitempty" json:"daemon"`
	Manager Manager `yaml:"manager,omitempty" json:"manager"`
	Log     Log     `yaml:"log,omitempty" json:"log"`
}

func initCmd() error {
	flag.StringVar(&file, "c", defaultFile, "configuration file")
	flag.BoolVar(&h, "h", false, "help")
	flag.BoolVar(&showFingerprint, "fingerprint", false, "show device fingerprint and exit")
	flag.Parse()
	if h {
		flag.Usage()
		return fmt.Errorf("invalid usage for command line")
	}
	return nil
}

// ShouldShowFingerprint 返回是否需要显示指纹
func ShouldShowFingerprint() bool {
	return showFingerprint
}

func Init() error {
	time.LoadLocation("Asia/Shanghai")

	err := initCmd()
	if err != nil {
		return err
	}

	err = initConf()
	if err != nil {
		return err
	}

	err = initLog()
	if err != nil {
		return err
	}

	return nil
}

func initConf() error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	Conf = &Configuration{}
	err = yaml.Unmarshal([]byte(data), Conf)
	return err
}

func initLog() error {
	level, err := log.ParseLevel(Conf.Log.Level)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	RotateLog = &lumberjack.Logger{
		Filename:   Conf.Log.File,
		MaxSize:    Conf.Log.MaxSize,
		MaxBackups: Conf.Log.MaxRolls,
	}
	log.SetOutput(RotateLog)
	return nil
}
