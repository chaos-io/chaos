package s3

import (
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/chaos-io/chaos/storage"
)

type mockS3API struct {
	s3iface.S3API
	putObjectWithContextFunc func(aws.Context, *awss3.PutObjectInput, ...request.Option) (*awss3.PutObjectOutput, error)
}

func (m *mockS3API) PutObjectWithContext(ctx aws.Context, input *awss3.PutObjectInput, opts ...request.Option) (*awss3.PutObjectOutput, error) {
	return m.putObjectWithContextFunc(ctx, input, opts...)
}

func TestS3ClientWrite(t *testing.T) {
	t.Run("writes object with default bucket", func(t *testing.T) {
		var gotBucket string
		var gotKey string
		var gotContentType string
		var gotContentLength int64
		var gotBody []byte

		client := &S3Client{
			s3: &mockS3API{
				putObjectWithContextFunc: func(_ aws.Context, input *awss3.PutObjectInput, _ ...request.Option) (*awss3.PutObjectOutput, error) {
					var err error
					gotBucket = aws.StringValue(input.Bucket)
					gotKey = aws.StringValue(input.Key)
					gotContentType = aws.StringValue(input.ContentType)
					gotContentLength = aws.Int64Value(input.ContentLength)
					gotBody, err = io.ReadAll(input.Body)
					if err != nil {
						t.Fatalf("unexpected body read error: %v", err)
					}
					return &awss3.PutObjectOutput{}, nil
				},
			},
			cfg: storage.NewConfig(func(cfg *storage.Config) {
				cfg.BucketName = "default-bucket"
			}),
		}

		obj := &storage.Object{
			Key:         "demo.txt",
			ContentType: "text/plain",
			Content:     []byte("hello"),
		}
		if err := client.Write(context.Background(), obj); err != nil {
			t.Fatalf("Write() error = %v", err)
		}

		if gotBucket != "default-bucket" {
			t.Fatalf("bucket = %q, want %q", gotBucket, "default-bucket")
		}
		if gotKey != "demo.txt" {
			t.Fatalf("key = %q, want %q", gotKey, "demo.txt")
		}
		if gotContentType != "text/plain" {
			t.Fatalf("contentType = %q, want %q", gotContentType, "text/plain")
		}
		if gotContentLength != int64(len("hello")) {
			t.Fatalf("contentLength = %d, want %d", gotContentLength, len("hello"))
		}
		if string(gotBody) != "hello" {
			t.Fatalf("body = %q, want %q", string(gotBody), "hello")
		}
	})

	t.Run("uses option bucket when provided", func(t *testing.T) {
		var gotBucket string

		client := &S3Client{
			s3: &mockS3API{
				putObjectWithContextFunc: func(_ aws.Context, input *awss3.PutObjectInput, _ ...request.Option) (*awss3.PutObjectOutput, error) {
					gotBucket = aws.StringValue(input.Bucket)
					return &awss3.PutObjectOutput{}, nil
				},
			},
			cfg: storage.NewConfig(func(cfg *storage.Config) {
				cfg.BucketName = "default-bucket"
			}),
		}

		err := client.Write(context.Background(), &storage.Object{
			Key:     "demo.txt",
			Content: []byte("hello"),
		}, storage.WithSignBucket("override-bucket"))
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}
		if gotBucket != "override-bucket" {
			t.Fatalf("bucket = %q, want %q", gotBucket, "override-bucket")
		}
	})

	t.Run("rejects mismatched size", func(t *testing.T) {
		client := &S3Client{
			cfg: storage.NewConfig(func(cfg *storage.Config) {
				cfg.BucketName = "default-bucket"
			}),
		}

		err := client.Write(context.Background(), &storage.Object{
			Key:     "demo.txt",
			Size:    10,
			Content: []byte("hello"),
		})
		if err == nil {
			t.Fatal("Write() error = nil, want non-nil")
		}
	})
}
