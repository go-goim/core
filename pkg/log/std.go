package log

import (
	"fmt"
	"log"
	"strings"
)

type stdLogger struct {
	logger *log.Logger
	option *option
}

func NewStdLogger(opts ...Option) Logger {
	options := newOption()
	for _, o := range opts {
		o(options)
	}

	s := &stdLogger{
		option: options,
	}
	s.logger = log.New(log.Writer(), options.prefix, log.LstdFlags)
	return s
}

func (s *stdLogger) Log(level Level, msg string, keyvals ...interface{}) {
	if level > s.option.level {
		return
	}

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "UNPAIRED_KEY")
	}

	s.logger.Printf("[%s] %s | %s", level.String(), msg, kvs2String(keyvals...))
}

func kvs2String(kvs ...interface{}) string {
	sb := strings.Builder{}

	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			continue
		}

		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(fmt.Sprintf("%v", kvs[i+1]))
		sb.WriteString(" ")
	}

	return sb.String()
}
