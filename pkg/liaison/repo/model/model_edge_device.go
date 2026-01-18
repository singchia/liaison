package model

import (
	"gorm.io/gorm"
)

const (
	TableNameEdgeDevice = "edge_device"
)

// EdgeDeviceRelationType Edge 和 Device 的关系类型
type EdgeDeviceRelationType int

const (
	// EdgeDeviceRelationHost Edge 所在 Device（Edge 运行在这个 Device 上）
	EdgeDeviceRelationHost EdgeDeviceRelationType = 1
	// EdgeDeviceRelationDiscovered Edge 发现的 Device（Edge 扫描发现的 Device）
	EdgeDeviceRelationDiscovered EdgeDeviceRelationType = 2
)

// EdgeDevice Edge 和 Device 的关系表
type EdgeDevice struct {
	gorm.Model
	EdgeID   uint64                `gorm:"column:edge_id;type:int;not null;uniqueIndex:idx_edge_device"`   // Edge ID
	DeviceID uint                   `gorm:"column:device_id;type:int;not null;uniqueIndex:idx_edge_device"` // Device ID
	Type     EdgeDeviceRelationType `gorm:"column:type;type:int;not null;uniqueIndex:idx_edge_device"`     // 关系类型：1. host, 2. discovered
}

func (EdgeDevice) TableName() string {
	return "edge_devices"
}
