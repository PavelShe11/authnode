package helpers

import (
	"github.com/PavelShe11/authnode/common/logger"
)

// NoopLogger implements logger.Logger interface with no-op methods for testing
type NoopLogger struct{}

func NewNoopLogger() logger.Logger {
	return &NoopLogger{}
}

func (l *NoopLogger) Debug(args ...interface{})                 {}
func (l *NoopLogger) Debugf(format string, args ...interface{}) {}
func (l *NoopLogger) Info(args ...interface{})                  {}
func (l *NoopLogger) Infof(format string, args ...interface{})  {}
func (l *NoopLogger) Warn(args ...interface{})                  {}
func (l *NoopLogger) Warnf(format string, args ...interface{})  {}
func (l *NoopLogger) Error(args ...interface{})                 {}
func (l *NoopLogger) Errorf(format string, args ...interface{}) {}
func (l *NoopLogger) Fatal(args ...interface{})                 {}
func (l *NoopLogger) Fatalf(format string, args ...interface{}) {}
