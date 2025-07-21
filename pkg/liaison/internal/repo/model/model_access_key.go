package model

import "gorm.io/gorm"

type AccessKey struct {
	gorm.Model
	EdgeID    uint   `gorm:"column:edge_id;type:int(11);not null"` // 关联的edge id
	AccessKey string `gorm:"column:access_key;type:varchar(255);not null"`
	SecretKey string `gorm:"column:secret_key;type:varchar(255);not null"`
}
