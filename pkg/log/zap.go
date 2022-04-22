package log

import (
	"path/filepath"
	"time"

	configv1 "github.com/yusank/goim/api/config/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	logger *zap.Logger
	option *option
}

func NewZapLogger(opts ...Option) Logger {
	options := newOption()
	for _, o := range opts {
		o(options)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(options.config.EncoderConfig),
		zapcore.AddSync(getLogWriter(options.outputPath)),
		zapcore.Level(int8(options.level-1)),
	)

	return &zapLogger{
		logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(4)),
		option: options,
	}
}

func getLogWriter(outputPath string) zapcore.WriteSyncer {
	// fileName is log file name contains current date
	fileName := getCurrentDate() + ".log"
	if outputPath != "" {
		fileName = filepath.Join(outputPath, fileName)
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

	switch level {
	case configv1.Level_DEBUG:
		z.logger.Debug(msg, kv2ZapFields(kvs...)...)
	case configv1.Level_INFO:
		z.logger.Info(msg, kv2ZapFields(kvs...)...)
	case configv1.Level_WARING:
		z.logger.Warn(msg, kv2ZapFields(kvs...)...)
	case configv1.Level_ERROR:
		z.logger.Error(msg, kv2ZapFields(kvs...)...)
	case configv1.Level_FATAL:
		z.logger.Fatal(msg, kv2ZapFields(kvs...)...)
	}
}

func kv2ZapFields(kvs ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			continue
		}

		fields = append(fields, zap.Any(key, kvs[i+1]))
	}
	return fields
}
