package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	configv1 "github.com/yusank/goim/api/config/v1"
)

type option struct {
	prefix      string
	level       configv1.Level
	config      zap.Config
	outputPath  string
	callerDepth int
}

func newOption() *option {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:   false,
		Encoding:      "json",
		EncoderConfig: encoderConfig,
	}
	return &option{
		outputPath:  "./logs",
		prefix:      "",
		level:       configv1.Level_DEBUG,
		config:      config,
		callerDepth: 4,
	}
}

type Option func(*option)

func WithZapConfig(config zap.Config) Option {
	return func(option *option) {
		option.config = config
	}
}

func WithPrefix(prefix string) Option {
	return func(option *option) {
		option.prefix = prefix
	}
}

func WithLevel(level configv1.Level) Option {
	return func(option *option) {
		option.level = level
	}
}

func WithOutputPath(path string) Option {
	return func(option *option) {
		option.outputPath = path
	}
}

func WithCallerDepth(d int) Option {
	return func(o *option) {
		o.callerDepth = d
	}
}
