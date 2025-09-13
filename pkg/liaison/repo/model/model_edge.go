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

type Edge struct {
	gorm.Model
	Name        string           `gorm:"column:name;type:varchar(255);not null"`
	Online      EdgeOnlineStatus `gorm:"column:online;type:int;not null"`
	HeartbeatAt time.Time        `gorm:"column:heartbeat_at;type:datetime;not null"`
	Description string           `gorm:"column:description;type:varchar(255);not null"`
	DeviceID    uint             `gorm:"column:device_id;type:int;not null"` // 关联设备
}
