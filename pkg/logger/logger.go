// Package logger provides initialization of global zap.Logger and simple access to it.
package logger

import (
	"go.uber.org/zap"
)

// Field is a key-value pair for providing context to logs.
type Field struct {
	Key   string
	Value string
}

// Logger is an entity for logging any data.
type Logger interface {
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Warning(message string, fields ...Field)
	Error(message string, fields ...Field)
}

// --- ZAP LOGGER ---

type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates new logger.Logger with zap.Logger underneath.
func NewZapLogger(logger *zap.Logger) Logger {
	return &zapLogger{
		logger: logger,
	}
}

// Debug logs messages with DEBUG importance level.
func (l *zapLogger) Debug(message string, fields ...Field) {
	l.logger.Debug(message, l.fieldsToZap(fields)...)
}

// Info logs messages with INFO importance level.
func (l *zapLogger) Info(message string, fields ...Field) {
	l.logger.Info(message, l.fieldsToZap(fields)...)
}

// Warning logs messages with WARNING importance level.
func (l *zapLogger) Warning(message string, fields ...Field) {
	l.logger.Warn(message, l.fieldsToZap(fields)...)
}

// Error logs messages with ERROR importance level.
func (l *zapLogger) Error(message string, fields ...Field) {
	l.logger.Error(message, l.fieldsToZap(fields)...)
}

func (l *zapLogger) fieldsToZap(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.String(f.Key, f.Value)
	}

	return zapFields
}

// --- NOP LOGGER ---

type nopLogger struct {
}

// NewNopLogger creates new logger.Logger that does nothing.
func NewNopLogger() Logger {
	return &nopLogger{}
}

// Debug does nothing but pretends it's doing something.
func (l *nopLogger) Debug(message string, fields ...Field) {
	return //nolint:all need for test coverage
}

// Info does nothing but pretends it's doing something.
func (l *nopLogger) Info(message string, fields ...Field) {
	return //nolint:all need for test coverage
}

// Warning does nothing but pretends it's doing something.
func (l *nopLogger) Warning(message string, fields ...Field) {
	return //nolint:all need for test coverage
}

// Error does nothing but pretends it's doing something.
func (l *nopLogger) Error(message string, fields ...Field) {
	return //nolint:all need for test coverage
}
