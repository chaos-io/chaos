package ck

import (
	"context"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"

	ckgo "github.com/ClickHouse/clickhouse-go/v2"
)

type Provider interface {
	NewSession(ctx context.Context) *gorm.DB
}

type provider struct {
	db *gorm.DB
}

func (p *provider) NewSession(ctx context.Context) *gorm.DB {
	return p.db.WithContext(ctx)
}

func NewCKFromConfig(cfg *Config) (Provider, error) {
	opt := &ckgo.Options{
		Addr: []string{cfg.Host},
		Auth: ckgo.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Debug:       cfg.Debug,
		DialTimeout: cfg.DialTimeout,
		HttpHeaders: cfg.HttpHeaders,
		ReadTimeout: cfg.ReadTimeout,
	}

	switch cfg.CompressionMethod {
	case CompressionMethodLZ4:
		opt.Compression = &ckgo.Compression{
			Method: ckgo.CompressionLZ4,
			Level:  cfg.CompressionLevel,
		}
	case CompressionMethodZSTD:
		opt.Compression = &ckgo.Compression{
			Method: ckgo.CompressionZSTD,
			Level:  cfg.CompressionLevel,
		}
	case CompressionMethodGZIP:
		opt.Compression = &ckgo.Compression{
			Method: ckgo.CompressionGZIP,
			Level:  cfg.CompressionLevel,
		}
	case CompressionMethodDeflate:
		opt.Compression = &ckgo.Compression{
			Method: ckgo.CompressionDeflate,
			Level:  cfg.CompressionLevel,
		}
	case CompressionMethodBrotli:
		opt.Compression = &ckgo.Compression{
			Method: ckgo.CompressionBrotli,
			Level:  cfg.CompressionLevel,
		}
	}

	switch cfg.Protocol {
	case ProtocolHTTP:
		opt.Protocol = ckgo.HTTP
	case ProtocolNative:
		opt.Protocol = ckgo.Native
	}

	ckSqlDB := ckgo.OpenDB(opt)
	db, err := gorm.Open(clickhouse.New(clickhouse.Config{Conn: ckSqlDB}))
	if err != nil {
		return nil, err
	}

	return &provider{db: db}, nil
}
