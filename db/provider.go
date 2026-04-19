package db

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type Provider interface {
	Gorm() *gorm.DB
	SQLDB() (*sql.DB, error)
	Close(ctx context.Context) error
}

var _ Provider = (*DB)(nil)

func (d *DB) Gorm() *gorm.DB {
	if d == nil {
		return nil
	}
	return d.DB
}

func (d *DB) SQLDB() (*sql.DB, error) {
	if d == nil || d.DB == nil {
		return nil, nil
	}
	return d.DB.DB()
}

func (d *DB) Close(ctx context.Context) error {
	_ = ctx

	sqlDB, err := d.SQLDB()
	if err != nil || sqlDB == nil {
		return err
	}
	return sqlDB.Close()
}
