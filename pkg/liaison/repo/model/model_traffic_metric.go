package model

import (
	"time"

	"gorm.io/gorm"
)

// TrafficMetric 流量监控数据模型
// 每分钟记录一次应用的流量统计
type TrafficMetric struct {
	gorm.Model
	ApplicationID uint      `gorm:"column:application_id;type:int;not null;index"` // 应用ID
	ProxyID       uint      `gorm:"column:proxy_id;type:int;not null;index"`       // 代理ID
	Timestamp     time.Time `gorm:"column:timestamp;type:datetime;not null;index"` // 时间戳（分钟级别）
	BytesIn       int64     `gorm:"column:bytes_in;type:bigint;not null;default:0"` // 入站流量（字节）
	BytesOut      int64     `gorm:"column:bytes_out;type:bigint;not null;default:0"` // 出站流量（字节）
	// 以下用于中间使用
	Application *Application `gorm:"-"`
	Proxy       *Proxy       `gorm:"-"`
}

func (TrafficMetric) TableName() string {
	return "traffic_metrics"
}
