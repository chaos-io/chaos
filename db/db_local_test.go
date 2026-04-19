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
		// DSN:    "root:@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=1s&readTimeout=3s&writeTimeout=3s",
		Debug: true,
	}

	if db, err := NewWithConfig(cfg); err != nil || db == nil {
		t.Errorf("NewWithConfig() = (%v, %v), want non-nil db and nil error", db, err)
	}
}

func Test_SqliteNew(t *testing.T) {
	cfg := &Config{
		Driver: SqliteDriver,
		DSN:    "test.db",
		Debug:  true,
	}

	if db, err := NewWithConfig(cfg); err != nil || db == nil {
		t.Errorf("NewWithConfig() = (%v, %v), want non-nil db and nil error", db, err)
	}
}

func Test_PgsqlNew(t *testing.T) {
	cfg := &Config{
		Driver: PostgresDriver,
		DSN:    "host=localhost user=eric dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		Debug:  true,
	}

	if db, err := NewWithConfig(cfg); err != nil || db == nil {
		t.Errorf("NewWithConfig() = (%v, %v), want non-nil db and nil error", db, err)
	}
}
