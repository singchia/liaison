package model

import "gorm.io/gorm"

type AccessKey struct {
	gorm.Model
	AccessKey string `gorm:"column:access_key;type:varchar(255);not null"`
	SecretKey string `gorm:"column:secret_key;type:varchar(255);not null"`
}
