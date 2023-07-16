package db

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	DriverName         string
	DSN                string
	MaxOpenConns       int           // default: 12
	MaxIdleConns       int           // default: 12
	ConnMaxLifetime    time.Duration // default: 2h
	ConnMaxIdleTime    time.Duration
	Debug              bool
	LogMode            bool
	DisableAutoMigrate bool

	*gorm.Config
}
