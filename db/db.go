package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func GetDB(logLevel logger.LogLevel) (*gorm.DB, error) {
	if db == nil {
		if err := initDB(logLevel); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// 初始化数据库
func initDB(logLevel logger.LogLevel) error {
	var err error
	db, err = gorm.Open(sqlite.Open("./gorm.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}
