package model

import (
	"time"

	"gorm.io/gorm"
)

type Edge struct {
	gorm.Model
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Description string    `gorm:"column:description;type:varchar(255);not null"`
	DeviceID    uint      `gorm:"column:device_id;type:int;not null"` // 关联设备
	AccessKeyID uint      `gorm:"column:access_key_id;type:int;not null"`
	HeartbeatAt time.Time `gorm:"column:heartbeat_at;type:datetime;not null"`
}
