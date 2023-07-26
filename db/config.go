package db

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Driver             string        `json:"driver"`
	DSN                string        `json:"dsn"`
	MaxOpenConns       int           `json:"maxOpenConns" default:"12"` // default: 12
	MaxIdleConns       int           `json:"maxIdleConns" default:"12"` // default: 12
	ConnMaxLifetime    time.Duration `json:"connMaxLifetime" `          // default: 2h
	ConnMaxIdleTime    time.Duration `json:"connMaxIdleTime"`
	Debug              bool          `json:"debug"`
	DisableAutoMigrate bool          `json:"disableAutoMigrate"`

	*gorm.Config
}
