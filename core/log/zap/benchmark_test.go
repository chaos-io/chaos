package zap

import (
	"errors"
	"testing"
	"time"
	
	log2 "github.com/chaos-io/chaos/core/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func BenchmarkZapLogger(b *testing.B) {
	// use config for both loggers
	cfg := NewDeployConfig()
	cfg.OutputPaths = nil
	cfg.ErrorOutputPaths = nil

	b.Run("stock", func(b *testing.B) {
		for _, level := range log2.Levels() {
			b.Run(level.String(), func(b *testing.B) {
				cfg.Level = zap.NewAtomicLevelAt(ZapifyLevel(level))

				logger, err := cfg.Build()
				require.NoError(b, err)

				funcs := []func(string, ...zap.Field){
					logger.Debug,
					logger.Info,
					logger.Warn,
					logger.Error,
					logger.Fatal,
				}

				message := "test"
				fields := []zap.Field{
					zap.String("test", "test"),
					zap.Bool("test", true),
					zap.Int("test", 42),
				}

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					funcs[i%(len(funcs)-1)](message, fields...)
				}
			})
		}
	})

	b.Run("wrapped", func(b *testing.B) {
		for _, level := range log2.Levels() {
			b.Run(level.String(), func(b *testing.B) {
				cfg.Level = zap.NewAtomicLevelAt(ZapifyLevel(level))
				logger, err := New(cfg)
				require.NoError(b, err)

				funcs := []func(string, ...log2.Field){
					logger.Debug,
					logger.Info,
					logger.Warn,
					logger.Error,
					logger.Fatal,
				}

				message := "test"
				fields := []log2.Field{
					log2.String("test", "test"),
					log2.Bool("test", true),
					log2.Int("test", 42),
				}

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					funcs[i%(len(funcs)-1)](message, fields...)
				}
			})
		}
	})
}

func BenchmarkZapifyField(b *testing.B) {
	fields := []log2.Field{
		log2.Nil("test"),
		log2.String("test", "test"),
		log2.Binary("test", []byte("test")),
		log2.Bool("test", true),
		log2.Int("test", 42),
		log2.UInt("test", 42),
		log2.Float64("test", 42),
		log2.Time("test", time.Now()),
		log2.Duration("test", time.Second),
		log2.NamedError("test", errors.New("test")),
		log2.Strings("test", []string{"test"}),
		log2.Any("test", "test"),
		log2.Reflect("test", "test"),
		log2.ByteString("test", []byte("test")),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zapifyField(fields[i%(len(fields)-1)])
	}
}

func BenchmarkZapifyFields(b *testing.B) {
	fields := []log2.Field{
		log2.Nil("test"),
		log2.String("test", "test"),
		log2.Binary("test", []byte("test")),
		log2.Bool("test", true),
		log2.Int("test", 42),
		log2.UInt("test", 42),
		log2.Float64("test", 42),
		log2.Time("test", time.Now()),
		log2.Duration("test", time.Second),
		log2.NamedError("test", errors.New("test")),
		log2.Strings("test", []string{"test"}),
		log2.Any("test", "test"),
		log2.Reflect("test", "test"),
		log2.ByteString("test", []byte("test")),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zapifyFields(fields...)
	}
}
