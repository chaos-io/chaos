package s3

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/chaos-io/chaos/infra/storage"
	"github.com/chaos-io/chaos/pkg/logs"
	"github.com/chaos-io/core/go/chaos/core"
	"github.com/samber/lo"
)

type S3Client struct {
	s3  *s3.S3
	cfg *storage.Config
}

func NewS3Client(cfg *storage.Config) (*S3Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	newSession, err := session.NewSession(&aws.Config{
		Region:           lo.ToPtr(cfg.Region),
		Endpoint:         lo.ToPtr(cfg.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		S3ForcePathStyle: lo.ToPtr(true),
	})
	if err != nil {
		return nil, logs.NewErrorw("failed to create S3 session", "error", err)
	}

	cli := s3.New(newSession)
	return &S3Client{s3: cli, cfg: cfg}, nil
}

func (c *S3Client) BucketName() string {
	return c.cfg.BucketName
}

func (c *S3Client) SetBucket(name string) error {
	return nil
}

func (c *S3Client) Read(key string, options core.Options) (*storage.Object, error) {
	ctx := context.Background()
	info, err := c.Stat(ctx, key, options)
	if err != nil {
		return nil, err
	}

	bucket := c.cfg.BucketName

	// small object, read in memory
	if info.Size < c.cfg.CacheSizeGT {
		input := &s3.GetObjectInput{
			Bucket: lo.ToPtr(bucket),
			Key:    lo.ToPtr(key),
		}
		obj, err := c.s3.GetObjectWithContext(ctx, input)
		if err != nil {
			return nil, logs.NewErrorf("failed to get object %s: %v", key, err)
		}

		data, err := io.ReadAll(obj.Body)
		if err != nil {
			return nil, logs.NewErrorf("failed to read object body(%s), error: %v", key, err)
		}

		return &storage.Object{
			LastModified: lo.FromPtr(obj.LastModified),
			Etag:         lo.FromPtr(obj.ETag),
			Key:          key,
			ContentType:  lo.FromPtr(obj.ContentType),
			Content:      data,
			Size:         lo.FromPtr(obj.ContentLength),
		}, nil
	}

	// big object, download to local file in parts
	tempFile, err := os.CreateTemp("", "s3-download-*.tmp")
	if err != nil {
		return nil, logs.NewErrorf("failed to create temp file, error: %v", err)
	}

	logs.Debugw("s3 object will be download, %s at %s, size=%d", key, tempFile.Name(), info.Size)

	if err := c.download(ctx, key, tempFile, options); err != nil {
		return nil, err
	}

	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return nil, logs.NewErrorf("failed to seek to temp file, error: %v", err)
	}

	data, err := io.ReadAll(tempFile)
	if err != nil {
		return nil, logs.NewErrorf("failed to read temp file, error: %v", err)
	}

	return &storage.Object{
		LastModified: info.LastModified,
		Key:          key,
		ContentType:  info.ContentType,
		Content:      data,
		Size:         info.Size,
	}, nil
}

func (c *S3Client) Write(obj *storage.Object, options core.Options) error {
	return nil
}

func (c *S3Client) Download(key, path string, options core.Options) error {
	options = WithConcurrencyOption(options)
	ctx := context.Background()
	file, err := os.Create(path)
	if err != nil {
		return logs.NewErrorw("failed to create file, error: %v", err)
	}

	return c.download(ctx, key, file, options)
}

func WithConcurrencyOption(options core.Options) core.Options {
	opt := core.NewOptions("concurrency", storage.DefaultConcurrency)
	opt.Merge(options)
	return opt
}

func (c *S3Client) Upload(localFile, key string, options core.Options) error {
	options = WithConcurrencyOption(options)
	ctx := context.Background()
	bucket := c.cfg.BucketName

	file, err := os.Open(localFile)
	if err != nil {
		return logs.NewErrorw("failed to open local file, error: %v", err)
	}

	uploader := s3manager.NewUploaderWithClient(c.s3, func(u *s3manager.Uploader) {
		u.PartSize = c.cfg.UploadPartSize
		if options != nil {
			u.Concurrency = options["concurrency"].(int)
		}
	})

	output, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: lo.ToPtr(bucket),
		Key:    lo.ToPtr(key),
		Body:   file,
	})
	if err != nil {
		return logs.NewErrorw("failed to upload file, error: %v", err)
	}

	logs.Infof("uploaded object %q, upload_id=%s", key, output.UploadID)
	return nil
}

func (c *S3Client) Stat(ctx context.Context, key string, options core.Options) (*storage.Object, error) {
	bucket := c.cfg.BucketName
	input := &s3.HeadObjectInput{
		Bucket: lo.ToPtr(bucket),
		Key:    lo.ToPtr(key),
	}
	output, err := c.s3.HeadObjectWithContext(ctx, input, nil)
	if err != nil {
		return nil, logs.NewErrorf("failed to get s3 head object(%s): %v", key, err)
	}

	return &storage.Object{
		LastModified: lo.FromPtr(output.LastModified),
		Key:          key,
		ContentType:  lo.FromPtr(output.ContentType),
		Size:         lo.FromPtr(output.ContentLength),
	}, nil
}

func (c *S3Client) download(ctx context.Context, key string, tempFile io.WriterAt, options core.Options) error {
	input := &s3.GetObjectInput{
		Bucket: lo.ToPtr(c.cfg.BucketName),
		Key:    lo.ToPtr(key),
	}

	downloader := s3manager.NewDownloaderWithClient(c.s3, func(d *s3manager.Downloader) {
		d.PartSize = c.cfg.DownloadPartSize
		if options != nil {
			d.Concurrency = options["concurrency"].(int)
		}
	})

	if _, err := downloader.DownloadWithContext(ctx, tempFile, input); err != nil {
		return logs.NewErrorw("failed to download object %s: %v", key, err)
	}

	return nil
}
