package dao

import (
	"github.com/singchia/liaison/pkg/liaison/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dao struct {
	db *gorm.DB

	// config
	config *config.Configuration
}

func NewDao(config *config.Configuration) (*dao, error) {
	db, err := gorm.Open(sqlite.Open(config.Manager.DB), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.Exec("PRAGMA synchronous = OFF;")
	sqlDB.Exec("PRAGMA journal_mode = DELETE;")
	sqlDB.Exec("PRAGMA cache_size = -2000;") // 2MB cache
	sqlDB.Exec("PRAGMA temp_store = MEMORY;")
	sqlDB.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	sqlDB.Exec("PRAGMA mmap_size = 268435456;") // 256MB memory map size
	sqlDB.SetMaxOpenConns(0)

	return &dao{
		db:     db,
		config: config,
	}, nil
}

func (dao *dao) Close() error {
	sqlDB, err := dao.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
