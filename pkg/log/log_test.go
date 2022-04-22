package log

import (
	"testing"
)

func TestStdLog(t *testing.T) {
	logger := NewStdLogger()
	SetLogger(logger)

	Debug("hello", "name", "world")
	Info("hello", "name", "world")
	Error("hello", "name", "world")
	Warn("hello", "name", "world")
	Fatal("hello", "name", "world")
}

func TestZapLog(t *testing.T) {
	logger := NewZapLogger()
	SetLogger(logger)

	Debug("hello", "name", "world")
	Info("hello", "name", "world")
	Error("hello", "name", "world")
	Warn("hello", "name", "world")
	Fatal("hello", "name", "world")
}
