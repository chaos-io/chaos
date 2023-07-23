package example

import (
	"reflect"
	"testing"

	"github.com/chaos-io/chaos/redis"
)

func TestSet(t *testing.T) {
	type args struct {
		redis *redis.Redis
		key   interface{}
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "set",
			args: args{
				redis: InitRedis(),
				key:   10,
				value: "ten",
			},
			want:    "OK",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Set(tt.args.redis, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		redis *redis.Redis
		key   string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "get",
			args: args{
				redis: InitRedis(),
				key:   "10",
			},
			want:    "ten",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.redis, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDel(t *testing.T) {
	type args struct {
		redis *redis.Redis
		key   interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "del",
			args: args{
				redis: InitRedis(),
				key:   10,
			},
			want:    int64(1),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Del(tt.args.redis, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Del() got(%T) = %v, want(%T) %v", got, got, tt.want, tt.want)
			}
		})
	}
}
