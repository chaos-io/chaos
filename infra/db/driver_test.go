//go:build local
// +build local

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MysqlNew(t *testing.T) {
	cfg := &Config{
		Driver: MysqlDriver,
		DSN:    "root:@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local",
		Debug:  true,
	}

	db, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func Test_SqliteNew(t *testing.T) {
	cfg := &Config{
		Driver: SqliteDriver,
		DSN:    "test.db",
		Debug:  true,
	}

	db, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func Test_PgsqlNew(t *testing.T) {
	cfg := &Config{
		Driver: PostgresDriver,
		DSN:    "host=localhost user=eric dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		Debug:  true,
	}

	db, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
