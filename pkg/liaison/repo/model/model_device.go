package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	TableNameDevice            = "device"
	TableNameEthernetInterface = "ethernet_interface"
)

type DeviceOnlineStatus int

const (
	DeviceOnlineStatusOnline  = 1
	DeviceOnlineStatusOffline = 2
)

type Device struct {
	gorm.Model
	Fingerprint string             `gorm:"column:fingerprint;type:varchar(255);not null"`
	Name        string             `gorm:"column:name;type:varchar(255);not null"`
	Description string             `gorm:"column:description;type:varchar(255);not null"`
	HostName    string             `gorm:"column:host_name;type:varchar(255);not null"`
	Online      DeviceOnlineStatus `gorm:"column:online;type:int;not null"`
	HeartbeatAt time.Time          `gorm:"column:heartbeat_at;type:datetime;not null"`
	// device info
	CPU        int                 `gorm:"column:cpu;type:int;not null"`
	Memory     int                 `gorm:"column:memory;type:int;not null"`
	Interfaces []EthernetInterface `gorm:"-"`
	//Disk      int    `gorm:"column:disk;type:int;not null"`
	OS        string `gorm:"column:os;type:varchar(255);not null"`
	OSVersion string `gorm:"column:os_version;type:varchar(255);not null"`
	// usage
	CPUUsage    float32 `gorm:"column:cpu_usage;type:float;not null"`
	MemoryUsage float32 `gorm:"column:memory_usage;type:float;not null"`
	DiskUsage   float32 `gorm:"column:disk_usage;type:float;not null"`
}

type EthernetInterface struct {
	gorm.Model
	DeviceID uint   `gorm:"column:device_id;type:int;not null"`
	Name     string `gorm:"column:name;type:varchar(255);not null"`
	MAC      string `gorm:"column:mac;type:varchar(255);not null"`
	IP       string `gorm:"column:ip;type:varchar(255);not null"`
	Netmask  string `gorm:"column:netmask;type:varchar(255);not null"`
}
