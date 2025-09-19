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
	file        string
	defaultFile string = "./edge.yaml"
)

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
	Manager Manager `yaml:"manager,omitempty" json:"manager"`
	Log     Log     `yaml:"log,omitempty" json:"log"`
}

func initCmd() error {
	flag.StringVar(&file, "c", defaultFile, "configuration file")
	flag.BoolVar(&h, "h", false, "help")
	flag.Parse()
	if h {
		flag.Usage()
		return fmt.Errorf("invalid usage for command line")
	}
	return nil
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
