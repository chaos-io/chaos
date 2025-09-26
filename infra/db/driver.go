package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	MysqlDriver    = "mysql"
	SqliteDriver   = "sqlite"
	PostgresDriver = "postgres"
)

func New(cfg *Config, opts ...gorm.Option) (*gorm.DB, error) {
	var err error
	var d *gorm.DB

	switch cfg.Driver {
	case MysqlDriver:
		d, err = gorm.Open(mysql.Open(cfg.DSN), opts...)
	case SqliteDriver:
		d, err = gorm.Open(sqlite.Open(cfg.DSN), opts...)
	case PostgresDriver:
		d, err = gorm.Open(postgres.Open(cfg.DSN), opts...)
	default:
		err = fmt.Errorf("database %q is not support", cfg.Driver)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open database, error: %v", err)
	}
	if d == nil {
		return nil, fmt.Errorf("failed to open database")
	}

	db, err := d.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect database, error: %v", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return d, nil
}
