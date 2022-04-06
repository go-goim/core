package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		zapcore.AddSync(options.writer),
		zapcore.Level(int8(options.level-1)),
	)

	return &zapLogger{
		logger: zap.New(core),
		option: options,
	}
}

func (z *zapLogger) Log(level Level, msg string, kvs ...interface{}) {
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "UNPAIRED_KEY")
	}

	switch level {
	case LevelDebug:
		z.logger.Debug(msg, kv2ZapFields(kvs...)...)
	case LevelInfo:
		z.logger.Info(msg, kv2ZapFields(kvs...)...)
	case LevelWarn:
		z.logger.Warn(msg, kv2ZapFields(kvs...)...)
	case LevelError:
		z.logger.Error(msg, kv2ZapFields(kvs...)...)
	case LevelFatal:
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
