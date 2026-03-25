package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/chaos-io/core/go/chaos/core"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"

	"github.com/chaos-io/chaos/logs"
	"github.com/chaos-io/chaos/storage"
)

func init() {
	storage.Register("minio", NewMinio)
}

type Minio struct {
	client     *minio.Client
	bucketName string
}

func NewMinio(cfg *storage.Config) storage.Storage {
	m := &Minio{}

	if client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	}); err != nil {
		logs.Errorw(fmt.Sprintf("failed to minio connect %s", cfg.Endpoint), "error", err)
		return nil
	} else {
		m.client = client
	}

	if err := m.SetBucket(context.Background(), cfg.BucketName); err != nil {
		logs.Errorw("failed to get bucket name", "error", err)
		return nil
	}

	return m
}

func (m *Minio) BucketName() string {
	return m.bucketName
}

func (m *Minio) SetBucket(ctx context.Context, name string) error {
	if len(name) > 0 && name != m.bucketName {
		if ok, err := m.client.BucketExists(ctx, name); err != nil {
			return err
		} else if !ok {
			if err = m.client.MakeBucket(ctx, name, minio.MakeBucketOptions{}); err != nil {
				return logs.NewErrorw("minio failed to create bucket", "bucketName", name, "error", err)
			}
		}
		m.bucketName = name
	}

	return nil
}

func (m *Minio) Read(ctx context.Context, key string, options core.Options) (*storage.Object, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		errResponse := &minio.ErrorResponse{}
		if errors.As(err, errResponse) && errResponse.Code == "NoSuchKey" && strings.HasPrefix(key, "/") {
			key = key[1:]
			bucketName := m.bucketName
			if pos := strings.Index(key, "/"); pos > 0 {
				bucketName = key[0:pos]
				key = key[pos:]
			}
			if obj, err = m.client.GetObject(ctx, bucketName, key, minio.GetObjectOptions{}); err != nil {
				if errors.As(err, errResponse) && errResponse.Code == "NoSuchKey" {
					return nil, logs.NewErrorw("failed to found the key", "key", key)
				}
				return nil, err
			}
		}
	}

	info, err := obj.Stat()
	if err != nil {
		errResponse := &minio.ErrorResponse{}
		if errors.As(err, errResponse) && errResponse.Code == "NoSuchKey" {
			return nil, logs.NewErrorw("failed to found the key", "key", key)
		}
		return nil, err
	}

	defer func() {
		_ = obj.Close()
	}()

	object := &storage.Object{
		Etag:         info.ETag,
		Key:          info.Key,
		LastModified: info.LastModified,
		Size:         info.Size,
		ContentType:  info.ContentType,
	}

	object.Content = make([]byte, info.Size)
	if size, err := obj.Read(object.Content); (err != nil && err != io.EOF) || size != int(info.Size) {
		return nil, logs.NewErrorw("failed to read the content from the minio object", "error", err)
	}

	return object, nil
}

func (m *Minio) Write(ctx context.Context, object *storage.Object, options core.Options) error {
	_, err := m.client.PutObject(ctx, m.bucketName, object.Key, bytes.NewReader(object.Content), object.Size, minio.PutObjectOptions{})
	return err
}

func (m *Minio) Download(ctx context.Context, key string, path string, options core.Options) error {
	if err := m.client.FGetObject(ctx, m.bucketName, key, path, minio.GetObjectOptions{}); err != nil {
		errResponse := &minio.ErrorResponse{}
		if errors.As(err, errResponse) && errResponse.Code == "NoSuchKey" && strings.HasPrefix(key, "/") {
			key = key[1:]
			bucketName := m.bucketName
			if pos := strings.Index(key, "/"); pos > 0 {
				bucketName = key[0:pos]
				key = key[pos:]
			}
			if err = m.client.FGetObject(ctx, bucketName, key, path, minio.GetObjectOptions{}); err == nil {
				return nil
			}
		}

		return logs.NewErrorw("minio failed to download object", "key", key, "path", path, "error", err)
	}

	return nil
}

func (m *Minio) Upload(ctx context.Context, localFile string, key string, options core.Options) error {
	if _, err := m.client.FPutObject(ctx, m.bucketName, key, localFile, minio.PutObjectOptions{}); err != nil {
		return logs.NewErrorw(fmt.Sprintf("minio failed to upload %s to %s", localFile, key), "error", err)
	}

	return nil
}
