package log

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	configv1 "github.com/go-goim/api/config/v1"
)

type zapLogger struct {
	logger *zap.Logger
	option *option
}

func NewZapLogger(opts ...Option) Logger {
	o := newOption()
	o.apply(opts...)

	var consoleCore zapcore.Core
	if o.enableConsole {
		consoleCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(o.getEncoderConfigForConsole()),
			zapcore.AddSync(colorable.NewColorableStdout()),
			zapcore.Level(int8(o.level-1)))

		if o.onlyConsole {
			return &zapLogger{
				logger: zap.New(consoleCore),
				option: o,
			}
		}
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(o.encoderConfig),
		zapcore.AddSync(getLogWriter(o)),
		zapcore.Level(int8(o.level-1)),
	)

	if o.enableConsole {
		core = zapcore.NewTee(core, consoleCore)
	}

	return &zapLogger{
		logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(o.callerDepth)),
		option: o,
	}
}

func newDefaultLogger() Logger {
	return NewZapLogger(EnableConsole(true), OnlyConsole(true))
}

func getLogWriter(o *option) zapcore.WriteSyncer {
	// fileName is log file name contains current date
	fileName := o.filenamePrefix + getCurrentDate() + ".log"
	if o.outputPath != "" {
		fileName = filepath.Join(o.outputPath, fileName)
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1024 * 1024 * 1024, // 1G
		MaxBackups: 5,
		MaxAge:     30,
	})
}

// getCurrentDate returns current date in format of YYYY-MM-DD
func getCurrentDate() string {
	return time.Now().Format("2006-01-02")
}
func (z *zapLogger) Log(level configv1.Level, msg string, kvs ...interface{}) {
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "UNPAIRED_KEY")
	}

	for k, v := range z.option.meta {
		kvs = append(kvs, k, v)
	}

	msg = strings.ReplaceAll(msg, "\n", " ")
	msg = strings.ReplaceAll(msg, "\r", " ")
	fields := kv2ZapFields(kvs...)
	switch level {
	case configv1.Level_DEBUG:
		z.logger.Debug(msg, fields...)
	case configv1.Level_INFO:
		z.logger.Info(msg, fields...)
	case configv1.Level_WARNING:
		z.logger.Warn(msg, fields...)
	case configv1.Level_ERROR:
		z.logger.Error(msg, fields...)
	case configv1.Level_FATAL:
		z.logger.Fatal(msg, fields...)
	}
}

func kv2ZapFields(kvs ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			continue
		}

		var valueStr string
		switch v := kvs[i+1].(type) {
		case string:
			valueStr = v
		case error:
			valueStr = v.Error()
		default:
			valueStr = fmt.Sprintf("%v", v)
		}

		fields = append(fields, zap.String(unescape(key), unescape(valueStr)))
	}
	return fields
}

func unescape(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")

	return s
}
