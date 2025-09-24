package ck

import "time"

type CompressionMethod string
type Protocol string

const (
	CompressionMethodLZ4     CompressionMethod = "lz4"
	CompressionMethodZSTD    CompressionMethod = "zstd"
	CompressionMethodGZIP    CompressionMethod = "gzip"
	CompressionMethodDeflate CompressionMethod = "deflate"
	CompressionMethodBrotli  CompressionMethod = "br"

	ProtocolHTTP   Protocol = "http"
	ProtocolNative Protocol = "native"
)

type Config struct {
	Host              string            `json:"host"`
	Database          string            `json:"database"`
	Username          string            `json:"username"`
	Password          string            `json:"password"`
	CompressionMethod CompressionMethod `json:"compressionMethod"`
	CompressionLevel  int               `json:"compressionLevel"`
	Protocol          Protocol          `json:"protocol"`
	DialTimeout       time.Duration     `json:"dialTimeout"`
	ReadTimeout       time.Duration     `json:"readTimeout"`
	Debug             bool              `json:"debug"`
	HttpHeaders       map[string]string `json:"httpHeaders"`
}
