package server

import (
	"github.com/alecthomas/log4go"
)

type Logger struct {
	*log4go.Logger
}

// Debugf logs messages at DEBUG level.
func (t Logger) Debugf(format string, args ...interface{}) {
	t.Debug(format, args...)
}

// Infof logs messages at INFO level.
func (t Logger) Infof(format string, args ...interface{}) {
	t.Info(format, args...)
}

// Warnf logs messages at WARN level.
func (t Logger) Warnf(format string, args ...interface{}) {
	t.Warn(format, args...)
}

// Errorf logs messages at ERROR level.
func (t Logger) Errorf(format string, args ...interface{}) {
	t.Error(format, args...)
}

// Fatalf logs messages at FATAL level.
func (t Logger) Fatalf(format string, args ...interface{}) {
	t.Critical(format, args...)
}
