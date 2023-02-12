package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	uzap "go.uber.org/zap"

	log2 "github.com/chaos-io/chaos/core/log"
	"github.com/chaos-io/chaos/core/log/nop"
	"github.com/chaos-io/chaos/core/log/zap"
)

var (
	msg        = "msg"
	msgfmt     = "%s %s"
	msgfmtargs = []interface{}{"hello", "world"}
	key        = "key"
	value      = "value"
	withKey    = "withKey"
	withValue  = "withValue"
)

var loggersToTest = []struct {
	name    string
	factory func(level log2.Level) (log2.Logger, error)
}{
	{
		name: "Zap",
		factory: func(level log2.Level) (log2.Logger, error) {
			cfg := zap.JSONConfig(level)
			// Disable output
			cfg.OutputPaths = []string{}
			cfg.ErrorOutputPaths = []string{}
			return zap.New(cfg)
		},
	},
	{
		name: "ZapNop",
		factory: func(level log2.Level) (log2.Logger, error) {
			return &zap.Logger{
				L: uzap.NewNop(),
			}, nil
		},
	},
	{
		name: "Nop",
		factory: func(level log2.Level) (log2.Logger, error) {
			return &nop.Logger{}, nil
		},
	},
}

func TestLoggers(t *testing.T) {
	for _, loggerInput := range loggersToTest {
		for _, level := range log2.Levels() {
			t.Run("Construct "+loggerInput.name+level.String(), func(t *testing.T) {
				logger, err := loggerInput.factory(level)
				require.NoError(t, err)
				require.NotNil(t, logger)

				lfmt := logger.Fmt()
				require.NotNil(t, lfmt)

				l := lfmt.Structured()
				require.NotNil(t, l)
				require.Equal(t, logger, l)
			})

			t.Run("With "+loggerInput.name+level.String(), func(t *testing.T) {
				logger, err := loggerInput.factory(level)
				require.NoError(t, err)
				require.NotNil(t, logger)

				withField := log2.String(withKey, withValue)
				loggerWith := log2.With(logger, withField)
				require.NotNil(t, loggerWith)
			})

			t.Run("AddCallerSkip "+loggerInput.name+level.String(), func(t *testing.T) {
				logger, err := loggerInput.factory(level)
				require.NoError(t, err)
				require.NotNil(t, logger)

				loggerCallerSkip := log2.AddCallerSkip(logger, 1)
				require.NotNil(t, loggerCallerSkip)
			})

			// TODO: validate log output
			t.Run("Logger "+loggerInput.name+level.String(), func(t *testing.T) {
				logger, err := loggerInput.factory(level)
				require.NoError(t, err)
				require.NotNil(t, logger)

				logger.Trace(msg, log2.String(key, value))
				logger.Debug(msg, log2.String(key, value))
				logger.Info(msg, log2.String(key, value))
				logger.Warn(msg, log2.String(key, value))
				logger.Error(msg, log2.String(key, value))
				// TODO: test fatal
			})

			// TODO: validate log output
			t.Run("LoggerFMT "+loggerInput.name+level.String(), func(t *testing.T) {
				logger, err := loggerInput.factory(level)
				require.NoError(t, err)
				require.NotNil(t, logger)

				lfmt := logger.Fmt()
				require.NotNil(t, lfmt)

				lfmt.Tracef(msgfmt, msgfmtargs...)
				lfmt.Debugf(msgfmt, msgfmtargs...)
				lfmt.Infof(msgfmt, msgfmtargs...)
				lfmt.Warnf(msgfmt, msgfmtargs...)
				lfmt.Errorf(msgfmt, msgfmtargs...)
				// TODO: test fatal
			})
		}
	}
}
