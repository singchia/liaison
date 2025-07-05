package model

import (
	"time"

	"gorm.io/gorm"
)

type Edge struct {
	gorm.Model
	DeviceID    uint      `gorm:"column:device_id;type:int;not null"` // 关联设备
	AccessKeyID uint      `gorm:"column:access_key_id;type:int;not null"`
	HeartbeatAt time.Time `gorm:"column:heartbeat_at;type:datetime;not null"`
}
