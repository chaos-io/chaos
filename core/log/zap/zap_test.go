package zap

import (
	"testing"
	
	log2 "github.com/chaos-io/chaos/core/log"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewQloudLogger(t *testing.T) {
	logger, err := NewQloudLogger(log2.DebugLevel)
	assert.NoError(t, err)

	core, logs := observer.New(zap.DebugLevel)

	logger.L = logger.L.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return addQloudContext(core).(zapcore.Core)
	}))

	expectedMessage := "test message"

	logger.Info(expectedMessage, log2.String("package", "zap"))
	assert.Equal(t, 1, logs.Len())

	loggedEntry := logs.AllUntimed()[0]
	assert.Equal(t, zap.InfoLevel, loggedEntry.Level)
	assert.Equal(t, expectedMessage, loggedEntry.Message)
	assert.Equal(t,
		map[string]interface{}{
			"@fields": map[string]interface{}{
				"package": "zap",
			},
		},
		loggedEntry.ContextMap(),
	)
}

func TestLogger_FormattedMethods(t *testing.T) {
	testCases := []struct {
		lvl          log2.Level
		expectLogged []zapcore.Entry
	}{
		{log2.TraceLevel, []zapcore.Entry{
			{Level: zap.DebugLevel, Message: "test at trace"},
			{Level: zap.DebugLevel, Message: "test at debug"},
			{Level: zap.InfoLevel, Message: "test at info"},
			{Level: zap.WarnLevel, Message: "test at warn"},
			{Level: zap.ErrorLevel, Message: "test at error"},
		}},
		{log2.DebugLevel, []zapcore.Entry{
			{Level: zap.DebugLevel, Message: "test at trace"},
			{Level: zap.DebugLevel, Message: "test at debug"},
			{Level: zap.InfoLevel, Message: "test at info"},
			{Level: zap.WarnLevel, Message: "test at warn"},
			{Level: zap.ErrorLevel, Message: "test at error"},
		}},
		{log2.InfoLevel, []zapcore.Entry{
			{Level: zap.InfoLevel, Message: "test at info"},
			{Level: zap.WarnLevel, Message: "test at warn"},
			{Level: zap.ErrorLevel, Message: "test at error"},
		}},
		{log2.WarnLevel, []zapcore.Entry{
			{Level: zap.WarnLevel, Message: "test at warn"},
			{Level: zap.ErrorLevel, Message: "test at error"},
		}},
		{log2.ErrorLevel, []zapcore.Entry{
			{Level: zap.ErrorLevel, Message: "test at error"},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.lvl.String(), func(t *testing.T) {
			logger, err := New(ConsoleConfig(tc.lvl))
			assert.NoError(t, err)

			core, logs := observer.New(ZapifyLevel(tc.lvl))

			logger.L = logger.L.WithOptions(zap.WrapCore(func(_ zapcore.Core) zapcore.Core {
				return core
			}))

			for _, lvl := range log2.Levels() {
				switch lvl {
				case log2.TraceLevel:
					logger.Tracef("test at %s", lvl.String())
				case log2.DebugLevel:
					logger.Debugf("test at %s", lvl.String())
				case log2.InfoLevel:
					logger.Infof("test at %s", lvl.String())
				case log2.WarnLevel:
					logger.Warnf("test at %s", lvl.String())
				case log2.ErrorLevel:
					logger.Errorf("test at %s", lvl.String())
				case log2.FatalLevel:
					// skipping fatal
				}
			}

			loggedEntries := logs.AllUntimed()

			assert.Equal(t, len(tc.expectLogged), logs.Len(), cmp.Diff(tc.expectLogged, loggedEntries))

			for i, le := range loggedEntries {
				assert.Equal(t, tc.expectLogged[i].Level, le.Level)
				assert.Equal(t, tc.expectLogged[i].Message, le.Message)
			}
		})
	}
}