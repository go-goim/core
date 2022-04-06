package log

// define log api and implement

type Logger interface {
	Log(level Level, msg string, keyvals ...interface{})
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
func (h *helper) Log(level Level, msg string, keyvals ...interface{}) {
	for _, l := range h.logger {
		l.Log(level, msg, keyvals...)
	}
}

// Debug logs a message at debug level.
func (h *helper) Debug(msg string, keyvals ...interface{}) {
	h.Log(LevelDebug, msg, keyvals...)
}

// Info logs a message at info level.
func (h *helper) Info(msg string, keyvals ...interface{}) {
	h.Log(LevelInfo, msg, keyvals...)
}

// Warn logs a message at warn level.
func (h *helper) Warn(msg string, keyvals ...interface{}) {
	h.Log(LevelWarn, msg, keyvals...)
}

// Error logs a message at error level.
func (h *helper) Error(msg string, keyvals ...interface{}) {
	h.Log(LevelError, msg, keyvals...)
}

// Fatal logs a message at fatal level.
func (h *helper) Fatal(msg string, keyvals ...interface{}) {
	h.Log(LevelFatal, msg, keyvals...)
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
