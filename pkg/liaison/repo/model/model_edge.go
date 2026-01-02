package model

import (
	"time"

	"gorm.io/gorm"
)

type EdgeOnlineStatus int

const (
	EdgeOnlineStatusOnline  EdgeOnlineStatus = 1
	EdgeOnlineStatusOffline EdgeOnlineStatus = 2
)

type EdgeStatus int

const (
	EdgeStatusRunning EdgeStatus = 1 // 运行中
	EdgeStatusStopped EdgeStatus = 2 // 停止
)

type Edge struct {
	gorm.Model
	Name        string           `gorm:"column:name;type:varchar(255);not null"`
	Status      EdgeStatus       `gorm:"column:status;type:int;not null;default:1"` // 运行状态：1. running, 2. stopped
	Online      EdgeOnlineStatus `gorm:"column:online;type:int;not null"`
	HeartbeatAt time.Time        `gorm:"column:heartbeat_at;type:datetime;not null"`
	Description string           `gorm:"column:description;type:varchar(255);not null"`
	DeviceID    uint             `gorm:"column:device_id;type:int;not null"` // 关联设备
}

func (Edge) TableName() string {
	return "edges"
}
