package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

const (
	MysqlDriverName    = "mysql"
	SqliteDriverName   = "sqlite"
	PostgresDriverName = "postgres"
)

type DB struct {
	*gorm.DB
	Config *Config
}

func New(cfg *Config) *DB {
	var err error
	var d *gorm.DB

	cfg.Config = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			// SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		// Logger: logger.Default.LogMode(logger.Info),
	}

	switch cfg.Driver {
	case MysqlDriverName:
		d, err = gorm.Open(mysql.Open(cfg.DSN), cfg.Config)
	case SqliteDriverName:
		d, err = gorm.Open(sqlite.Open(cfg.DSN), cfg.Config)
	case PostgresDriverName:
		d, err = gorm.Open(postgres.Open(cfg.DSN), cfg.Config)
	default:
		err = fmt.Errorf("database %q is not support", cfg.Driver)
	}
	if err != nil {
		panic(fmt.Errorf("failed to connect database, error: %v", err))
	}

	if cfg.Debug {
		d = d.Debug()
	}

	db, err := d.DB()
	if err != nil {
		panic(fmt.Sprintf("get db error: %v", err))
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &DB{d, cfg}
}
