//go:build local
// +build local

package example

import (
	"path"
	"reflect"
	"testing"

	"github.com/chaos-io/chaos/storage"
	_ "github.com/chaos-io/chaos/storage/minio"
	"github.com/chaos-io/core/go/chaos/core"
)

func TestWrite(t *testing.T) {
	type args struct {
		object  *storage.Object
		options core.Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write",
			args: args{
				object: &storage.Object{
					Etag:    "",
					Key:     path.Join("test", "write_test1"),
					Size:    int64(len("write test 1")),
					Content: []byte("write test 1"),
				},
				options: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.Write(tt.args.object, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpload(t *testing.T) {
	type args struct {
		localFile string
		key       string
		options   core.Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "upload",
			args: args{
				localFile: "./config/storage.yaml",
				key:       path.Join("test", "storage.yaml"),
				options:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.Upload(tt.args.localFile, tt.args.key, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		key     string
		options core.Options
	}
	tests := []struct {
		name        string
		args        args
		wantContent []byte
		wantErr     bool
	}{
		{
			name: "read",
			args: args{
				key:     path.Join("test", "write_test1"),
				options: nil,
			},
			wantContent: []byte("write test 1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.Read(tt.args.key, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Content, tt.wantContent) {
				t.Errorf("Read() got = %s, wantContent %s", got.Content, tt.wantContent)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	type args struct {
		key     string
		path    string
		options core.Options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "download",
			args: args{
				key:     path.Join("test", "write_test1"),
				path:    "./write_test1",
				options: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.Download(tt.args.key, tt.args.path, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
