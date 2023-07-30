package db

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		args args
		name string
	}{
		{
			name: MysqlDriver,
			args: args{
				cfg: &Config{
					Driver: MysqlDriver,
					DSN:    "root:@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local",
					Debug:  true,
				},
			},
		},
		{
			name: SqliteDriver,
			args: args{
				cfg: &Config{
					Driver: SqliteDriver,
					DSN:    "test.db",
					Debug:  true,
				},
			},
		},
		{
			name: PostgresDriver,
			args: args{
				cfg: &Config{
					Driver: PostgresDriver,
					DSN:    "host=localhost user=eric dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai",
					Debug:  true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("new db panic, %v", r)
				}
			}()

			if db := New(tt.args.cfg); db == nil {
				t.Errorf("New() is nil")
			}
		})
	}
}
