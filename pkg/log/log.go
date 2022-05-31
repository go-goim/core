package log

import (
	kraoslogger "github.com/go-kratos/kratos/v2/log"

	configv1 "github.com/go-goim/goim/api/config/v1"
)

// define log api and implement

type Logger interface {
	Log(level configv1.Level, msg string, keyvals ...interface{})
}

var (
	global = newDefaultLogger()
)

// Debug logs a message at debug level.
func Debug(msg string, keyvals ...interface{}) {
	global.Log(configv1.Level_DEBUG, msg, keyvals...)
}

// Info logs a message at info level.
func Info(msg string, keyvals ...interface{}) {
	global.Log(configv1.Level_INFO, msg, keyvals...)
}

// Warn logs a message at warn level.
func Warn(msg string, keyvals ...interface{}) {
	global.Log(configv1.Level_WARNING, msg, keyvals...)
}

// Error logs a message at error level.
func Error(msg string, keyvals ...interface{}) {
	global.Log(configv1.Level_ERROR, msg, keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(msg string, keyvals ...interface{}) {
	global.Log(configv1.Level_FATAL, msg, keyvals...)
}

func SetLogger(logger Logger) {
	global = logger
}

func GetLogger() Logger {
	return global
}

// Kratos logger here

func SetKratosLogger(logger Logger) {
	kraoslogger.SetLogger(logger2KratosLogger(logger))
}

func logger2KratosLogger(l Logger) kraoslogger.Logger {
	return &loggerConvert{logger: l}
}

type loggerConvert struct {
	logger Logger
}

func (l *loggerConvert) Log(level kraoslogger.Level, keyvals ...interface{}) error {
	if len(keyvals) < 2 {
		return nil
	}

	msg, ok := keyvals[1].(string)
	if !ok {
		return nil
	}

	l.logger.Log(configv1.Level(level)+1, msg, keyvals[2:]...)
	return nil
}
