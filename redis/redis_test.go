package redis

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name:    "rdb",
			args:    args{cfg: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !tt.wantErr {
				t.Errorf("New() = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}
