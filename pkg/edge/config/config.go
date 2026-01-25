package config

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/singchia/liaison/pkg/config"
	"github.com/singchia/liaison/pkg/utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

// ErrShowFingerprint 表示需要显示指纹并退出的错误
var ErrShowFingerprint = errors.New("show fingerprint and exit")

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
	if showFingerprint {
		showFingerprintDetails()
		return ErrShowFingerprint
	}
	return nil
}

// ShouldShowFingerprint 返回是否需要显示指纹
func ShouldShowFingerprint() bool {
	return showFingerprint
}

// showFingerprintDetails 显示详细的指纹信息
func showFingerprintDetails() {
	// 使用工具函数获取指纹
	fingerprint, err := utils.GetFingerprint()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting fingerprint: %v\n", err)
		os.Exit(1)
	}

	// 获取详细信息用于显示
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting network interfaces: %v\n", err)
		os.Exit(1)
	}
	macs := make([]string, 0)
	macList := []string{}
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 && iface.Flags&net.FlagLoopback == 0 {
			macStr := iface.HardwareAddr.String()
			macs = append(macs, macStr)
			macList = append(macList, fmt.Sprintf("%s (%s)", macStr, iface.Name))
		}
	}
	sort.Strings(macs)
	mac := strings.Join(macs, ",")

	// CPU 信息
	cpuID := ""
	cpuInfo, err := cpu.Info()
	if err != nil {
		hostname, _ := os.Hostname()
		cpuID = hostname
		fmt.Printf("CPU Info: Error getting CPU info, using hostname as fallback: %s\n", hostname)
	} else if len(cpuInfo) > 0 {
		cpuStrings := make([]string, 0, len(cpuInfo))
		for _, info := range cpuInfo {
			cpuStr := info.ModelName + info.VendorID + info.Family
			cpuStrings = append(cpuStrings, cpuStr)
			fmt.Printf("CPU Info[%d]: ModelName=%s, VendorID=%s, Family=%s\n", len(cpuStrings)-1, info.ModelName, info.VendorID, info.Family)
		}
		cpuMap := make(map[string]bool)
		uniqueCpus := make([]string, 0)
		for _, cpuStr := range cpuStrings {
			if !cpuMap[cpuStr] {
				cpuMap[cpuStr] = true
				uniqueCpus = append(uniqueCpus, cpuStr)
			}
		}
		sort.Strings(uniqueCpus)
		cpuID = strings.Join(uniqueCpus, ",")
	}

	// 磁盘序列号
	diskIDs := make([]string, 0)
	counts, err := disk.IOCounters()
	if err == nil && len(counts) > 0 {
		for _, stats := range counts {
			if stats.SerialNumber != "" {
				diskIDs = append(diskIDs, stats.SerialNumber)
				fmt.Printf("Disk Serial: %s (from %s)\n", stats.SerialNumber, stats.Name)
			}
		}
	}
	if len(diskIDs) == 0 {
		hostname, _ := os.Hostname()
		diskIDs = []string{hostname}
		fmt.Printf("Disk Serial: Error getting disk serial, using hostname as fallback: %s\n", hostname)
	}
	sort.Strings(diskIDs)
	diskID := strings.Join(diskIDs, ",")

	// 显示详细信息
	fmt.Println("\n=== Fingerprint Details ===")
	fmt.Printf("All MAC Addresses:\n")
	for _, m := range macList {
		fmt.Printf("  - %s\n", m)
	}
	fmt.Printf("Sorted MACs (used): %s\n", mac)
	fmt.Printf("CPU ID (sorted): %s\n", cpuID)
	fmt.Printf("Disk ID (sorted): %s\n", diskID)

	// 计算原始字符串
	raw := fmt.Sprintf("%s|%s|%s", mac, cpuID, diskID)
	fmt.Printf("\nRaw String: %s\n", raw)
	fmt.Printf("Fingerprint: %s\n", fingerprint)
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
