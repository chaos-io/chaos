//go:build local
// +build local

package db

import (
	"testing"
)

func Test_MysqlNew(t *testing.T) {
	cfg := &Config{
		Driver: MysqlDriver,
		DSN:    "root:@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local",
		Debug:  true,
	}

	if db := New(cfg); db == nil {
		t.Errorf("New() is nil")
	}
}

func Test_SqliteNew(t *testing.T) {
	cfg := &Config{
		Driver: SqliteDriver,
		DSN:    "test.db",
		Debug:  true,
	}

	if db := New(cfg); db == nil {
		t.Errorf("New() is nil")
	}
}

func Test_PgsqlNew(t *testing.T) {
	cfg := &Config{
		Driver: PostgresDriver,
		DSN:    "host=localhost user=eric dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		Debug:  true,
	}

	if db := New(cfg); db == nil {
		t.Errorf("New() is nil")
	}
}
