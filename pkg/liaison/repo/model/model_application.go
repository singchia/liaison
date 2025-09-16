package model

import (
	"time"

	"gorm.io/gorm"
)

type ApplicationType string

const (
	ApplicationTypeTCP ApplicationType = "tcp" // currently only support tcp
)

// Application is a software running on a device that edge-proxy can proxy to.
type Application struct {
	gorm.Model
	EdgeIDs         []uint          `gorm:"column:edge_ids;type:json;not null"` // 关联的edge id
	DeviceID        uint            `gorm:"column:device_id;type:int;not null"`
	Name            string          `gorm:"column:name;type:varchar(255);not null"`
	IP              string          `gorm:"column:ip;type:varchar(255);not null"`
	Port            int             `gorm:"column:port;type:int;not null"`
	HeartbeatAt     time.Time       `gorm:"column:heartbeat_at;type:datetime;not null"`
	ApplicationType ApplicationType `gorm:"column:application_type;type:varchar(255);not null"`
}
