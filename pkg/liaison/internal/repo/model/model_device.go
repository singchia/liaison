package model

import "gorm.io/gorm"

const (
	TableNameDevice            = "device"
	TableNameEthernetInterface = "ethernet_interface"
)

type DeviceOnlineStatus int

type Device struct {
	gorm.Model
	Fingerprint string `gorm:"column:fingerprint;type:varchar(255);not null"`
	Name        string `gorm:"column:name;type:varchar(255);not null"`
	HostName    string `gorm:"column:host_name;type:varchar(255);not null"`
	CPU         int    `gorm:"column:cpu;type:int;not null"`
	Memory      int    `gorm:"column:memory;type:int;not null"`
	Disk        int    `gorm:"column:disk;type:int;not null"`
	OS          string `gorm:"column:os;type:varchar(255);not null"`
	Version     string `gorm:"column:version;type:varchar(255);not null"`
	CPUUsage    int    `gorm:"column:cpu_usage;type:int;not null"`
	MemoryUsage int    `gorm:"column:memory_usage;type:int;not null"`
	DiskUsage   int    `gorm:"column:disk_usage;type:int;not null"`
	Description string `gorm:"column:description;type:varchar(255);not null"`
}

type EthernetInterface struct {
	gorm.Model
	DeviceID uint   `gorm:"column:device_id;type:int;not null"`
	Name     string `gorm:"column:name;type:varchar(255);not null"`
	MAC      string `gorm:"column:mac;type:varchar(255);not null"`
	IP       string `gorm:"column:ip;type:varchar(255);not null"`
	Netmask  string `gorm:"column:netmask;type:varchar(255);not null"`
	Gateway  string `gorm:"column:gateway;type:varchar(255);not null"`
}
