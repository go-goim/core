package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	configv1 "github.com/go-goim/core/api/config/v1"
)

type option struct {
	filenamePrefix string
	level          configv1.Level
	encoderConfig  zapcore.EncoderConfig
	outputPath     string
	callerDepth    int
	enableConsole  bool
	onlyConsole    bool
	meta           map[string]interface{}
}

func newOption() *option {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return &option{
		outputPath:     "./logs",
		filenamePrefix: "",
		level:          configv1.Level_DEBUG,
		encoderConfig:  encoderConfig,
	}
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *option) getEncoderConfigForConsole() zapcore.EncoderConfig {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return encoderConfig
}

type Option func(*option)

func EncodeConfig(config zapcore.EncoderConfig) Option {
	return func(option *option) {
		option.encoderConfig = config
	}
}

func FilenamePrefix(prefix string) Option {
	return func(option *option) {
		option.filenamePrefix = prefix
	}
}

func Level(level configv1.Level) Option {
	return func(option *option) {
		option.level = level
	}
}

func OutputPath(path string) Option {
	return func(option *option) {
		option.outputPath = path
	}
}

func CallerDepth(d int) Option {
	return func(o *option) {
		o.callerDepth = d
	}
}

func EnableConsole(enable bool) Option {
	return func(o *option) {
		o.enableConsole = enable
	}
}

func OnlyConsole(only bool) Option {
	return func(o *option) {
		o.onlyConsole = only
	}
}

func Meta(key string, value interface{}) Option {
	return func(o *option) {
		if o.meta == nil {
			o.meta = make(map[string]interface{})
		}
		o.meta[key] = value
	}
}
