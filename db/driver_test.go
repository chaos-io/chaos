package db

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "db1",
			args: args{
				cfg: &Config{
					Driver: MysqlDriverName,
					DSN:    "root:@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local",
					Debug:  true,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !tt.wantErr {
				t.Errorf("New() = %v, want %v, err %v", got, tt.wantErr, tt.err)
			}
		})
	}
}
