package log

import (
	"io"

	"go.uber.org/zap"
)

type option struct {
	prefix string
	level  Level
	config zap.Config
	writer io.Writer
}

func newOption() *option {
	return &option{
		prefix: "",
		level:  LevelDebug,
		config: zap.NewProductionConfig(),
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

func WithLevel(level Level) Option {
	return func(option *option) {
		option.level = level
	}
}

func WithWriter(writer io.Writer) Option {
	return func(option *option) {
		option.writer = writer
	}
}
