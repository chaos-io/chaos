package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

const (
	MysqlDriverName    = "mysql"
	SqliteDriverName   = "sqlite"
	PostgresDriverName = "postgres"
)

var DB *gorm.DB

func New(cfg *Config) *gorm.DB {
	var err error

	cfg.Config = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			// SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch cfg.DriverName {
	case MysqlDriverName:
		DB, err = gorm.Open(mysql.Open(cfg.DSN), cfg.Config)
	case SqliteDriverName:
		DB, err = gorm.Open(sqlite.Open(cfg.DSN), cfg.Config)
	case PostgresDriverName:
		DB, err = gorm.Open(postgres.Open(cfg.DSN), cfg.Config)
	default:
		err = fmt.Errorf("database %s is not support", cfg.DriverName)
	}
	if err != nil {
		panic("failed to connect database")
	}

	if cfg.Debug {
		DB = DB.Debug()
	}

	db, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("get db error: %v", err))
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return DB
}
