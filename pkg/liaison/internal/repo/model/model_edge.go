package model

import (
	"time"

	"gorm.io/gorm"
)

type EdgeStatus int

const (
	EdgeStatusOffline EdgeStatus = 1
	EdgeStatusRunning EdgeStatus = 2
	EdgeStatusStopped EdgeStatus = 3
)

type Edge struct {
	gorm.Model
	Name        string     `gorm:"column:name;type:varchar(255);not null"`
	Status      EdgeStatus `gorm:"column:status;type:int;not null"`
	Description string     `gorm:"column:description;type:varchar(255);not null"`
	DeviceID    uint       `gorm:"column:device_id;type:int;not null"` // 关联设备
	HeartbeatAt time.Time  `gorm:"column:heartbeat_at;type:datetime;not null"`
}
