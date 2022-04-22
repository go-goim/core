package log

import (
	kraoslogger "github.com/go-kratos/kratos/v2/log"

	configv1 "github.com/yusank/goim/api/config/v1"
)

// define log api and implement

type Logger interface {
	Log(level configv1.Level, msg string, keyvals ...interface{})
}

var global = newHelper(NewStdLogger())

// helper is a logger helper.
// helper implements leveled logger.
type helper struct {
	logger []Logger
}

func newHelper(logger ...Logger) *helper {
	return &helper{logger: logger}
}

// Log print log by level and keyvals.
func (h *helper) Log(level configv1.Level, msg string, keyvals ...interface{}) {
	for _, l := range h.logger {
		l.Log(level, msg, keyvals...)
	}
}

// Debug logs a message at debug level.
func (h *helper) Debug(msg string, keyvals ...interface{}) {
	h.Log(configv1.Level_DEBUG, msg, keyvals...)
}

// Info logs a message at info level.
func (h *helper) Info(msg string, keyvals ...interface{}) {
	h.Log(configv1.Level_INFO, msg, keyvals...)
}

// Warn logs a message at warn level.
func (h *helper) Warn(msg string, keyvals ...interface{}) {
	h.Log(configv1.Level_WARING, msg, keyvals...)
}

// Error logs a message at error level.
func (h *helper) Error(msg string, keyvals ...interface{}) {
	h.Log(configv1.Level_ERROR, msg, keyvals...)
}

// Fatal logs a message at fatal level.
func (h *helper) Fatal(msg string, keyvals ...interface{}) {
	h.Log(configv1.Level_FATAL, msg, keyvals...)
}

// Debug logs a message at debug level.
func Debug(msg string, keyvals ...interface{}) {
	global.Debug(msg, keyvals...)
}

// Info logs a message at info level.
func Info(msg string, keyvals ...interface{}) {
	global.Info(msg, keyvals...)
}

// Warn logs a message at warn level.
func Warn(msg string, keyvals ...interface{}) {
	global.Warn(msg, keyvals...)
}

// Error logs a message at error level.
func Error(msg string, keyvals ...interface{}) {
	global.Error(msg, keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(msg string, keyvals ...interface{}) {
	global.Fatal(msg, keyvals...)
}

func SetLogger(logger ...Logger) {
	global = newHelper(logger...)
}

func GetLogger() Logger {
	return global
}

func SetKratosLogger(logger Logger) {
	kraoslogger.SetLogger(logger2KratosLogger(logger))
}

func logger2KratosLogger(l Logger) kraoslogger.Logger {
	return &loggerConvert{l}
}

type loggerConvert struct {
	l Logger
}

func (l *loggerConvert) Log(level kraoslogger.Level, keyvals ...interface{}) error {
	if len(keyvals) < 2 {
		return nil
	}

	msg, ok := keyvals[1].(string)
	if !ok {
		return nil
	}

	l.l.Log(configv1.Level(level)+1, msg, keyvals[2:]...)
	return nil
}
