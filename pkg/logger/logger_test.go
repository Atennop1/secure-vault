package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewZapLogger(t *testing.T) {
	initial, _ := zap.NewProduction()
	logger := NewZapLogger(initial)
	require.Equal(t, logger.(*zapLogger).logger, initial)
}

func TestZapLogger(t *testing.T) {
	cases := []struct {
		name    string
		level   zapcore.Level
		message string
	}{
		{
			name:    "Debug",
			level:   zap.DebugLevel,
			message: "something",
		},
		{
			name:    "Info",
			level:   zap.InfoLevel,
			message: "something",
		},
		{
			name:    "Warning",
			level:   zap.WarnLevel,
			message: "something",
		},
		{
			name:    "Error",
			level:   zap.ErrorLevel,
			message: "something",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(tc.level)
			observedLogger := zap.New(observedZapCore)
			logger := NewZapLogger(observedLogger)

			switch tc.level {
			case zap.DebugLevel:
				logger.Debug(tc.message)
			case zap.InfoLevel:
				logger.Info(tc.message)
			case zap.WarnLevel:
				logger.Warning(tc.message)
			case zap.ErrorLevel:
				logger.Error(tc.message)
			}

			require.Equal(t, 1, observedLogs.Len())
			require.Equal(t, tc.level, observedLogs.All()[0].Level)
			require.Equal(t, tc.message, observedLogs.All()[0].Message)
		})
	}
}

func TestFieldsToZap(t *testing.T) {
	zapLogger := zapLogger{logger: zap.Must(zap.NewProduction())}
	require.Equal(t, []zap.Field{
		zap.String("key1", "value1"),
		zap.String("key2", "value2"),
	}, zapLogger.fieldsToZap([]Field{
		{"key1", "value1"},
		{"key2", "value2"},
	}))
}

func TestNewNopLogger(t *testing.T) {
	_ = NewNopLogger()
	require.True(t, true)
}

func TestNopLogger(t *testing.T) {
	cases := []struct {
		name  string
		level zapcore.Level
	}{
		{name: "Debug", level: zap.DebugLevel},
		{name: "Info", level: zap.InfoLevel},
		{name: "Warning", level: zap.WarnLevel},
		{name: "Error", level: zap.ErrorLevel},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger := NewNopLogger()

			switch tc.name {
			case "Debug":
				logger.Debug("nop")
			case "Info":
				logger.Info("nop")
			case "Warning":
				logger.Warning("nop")
			case "Error":
				logger.Error("nop")
			}
		})
	}
}
