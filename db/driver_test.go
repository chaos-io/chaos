package db

import (
	"log/slog"
	"maps"
	"slices"
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

func TestMin(t *testing.T) {
	slog.Info("min", "min", min("a", "b"))
	slog.Info("min", "min", min(1, 2))
	slog.Info("min", "min", min('a', 'b'))
}

func TestClear(t *testing.T) {
	sli := []string{"a", "b"}
	slog.Info("string slice", "sli", sli)
	clear(sli)
	slog.Info("string slice cleared", "sli", sli)

	m := map[string]int{"a": 1, "b": 2}
	slog.Info("map", "m", m)
	clear(m)
	slog.Info("map cleared", "m", m)

	var params = []struct {
		Id   string
		Name string
	}{{Id: "1", Name: "a"}}
	slog.Info("params", "params", params)
	clear(params)
	slog.Info("params", "params", params)
}

func TestSlices(t *testing.T) {
	sli := []string{"a", "b", "c"}
	slog.Info("slice first index", "", slices.Index(sli, "a"))
	slog.Info("slice contains", "", slices.Contains(sli, "a"))
	slog.Info("slice contains", "", slices.Clip(sli))
	slices.Reverse(sli)
	slog.Info("slice contains", "", sli)
	slog.Info("slice delete", "", slices.Delete(sli, 1, 2)) // 左闭右开
}

func TestMaps(t *testing.T) {
	m := map[string]interface{}{
		"a": 1,
		"b": "c",
	}
	slog.Info("map", "", m)

	clone := maps.Clone(m)
	slog.Info("map clone", "", clone)

	m2 := map[string]interface{}{}
	maps.Copy(m2, m)
	slog.Info("map copy", "", m2)

	slog.Info("map equal1", "", maps.Equal(m, clone))
	slog.Info("map equal2", "", maps.Equal(m, m2))

	slog.Info("map equalFunc", "", maps.EqualFunc(m, m2, func(i interface{}, i2 interface{}) bool {
		return m["a"] == m2["a"]
	}))

	maps.DeleteFunc(m, func(s string, i interface{}) bool {
		return s == "a"
	})
	slog.Info("map deleted", "", m)
}
