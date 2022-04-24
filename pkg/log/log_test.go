package log

import (
	"testing"
)

func TestZapLog(t *testing.T) {
	logger := NewZapLogger()
	SetLogger(logger)

	Debug("hello", "name", "world")
	Info("hello", "name", "world")
	Error("hello", "name", "world")
	Warn("hello", "name", "world")
}
