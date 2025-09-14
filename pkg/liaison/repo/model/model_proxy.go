package model

import "gorm.io/gorm"

type ProxyStatus int

const (
	ProxyStatusRunning ProxyStatus = iota + 1
	ProxyStatusStopped
)

const (
	TableNameProxy = "proxy"
)

type Proxy struct {
	gorm.Model
	ApplicationID uint        `gorm:"column:application_id;type:int;not null"`
	Name          string      `gorm:"column:name;type:varchar(255);not null"`
	Port          int         `gorm:"column:port;type:int;not null"`
	Status        ProxyStatus `gorm:"column:status;type:int;not null"`
	Description   string      `gorm:"column:description;type:varchar(255);not null"`
	// 以下用于中间使用
	Application *Application `gorm:"-"`
}
