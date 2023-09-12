package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

const (
	level = "level"
)

type Level int

type Logger interface {
	Printf(format string, v ...interface{})
	Level(l Level) Logger
	WithStr(k, v string) Logger
	WithInt(k string, v int) Logger
	WithFlo(k string, v float64) Logger
	WithBoo(k string, v bool) Logger
}

type ZeroLoggerWrapper struct {
	logger   zerolog.Logger
	LogLevel Level
}

func New(
	logger zerolog.Logger,
	options ...func(*ZeroLoggerWrapper),
) *ZeroLoggerWrapper {
	wrapper := &ZeroLoggerWrapper{
		logger:   logger,
		LogLevel: 0,
	}

	for _, option := range options {
		option(wrapper)
	}

	return wrapper
}

func (log *ZeroLoggerWrapper) Level(l Level) Logger {
	if l < log.LogLevel {
		return &NoopLogger{}
	}

	return &ZeroLoggerWrapper{
		logger: log.logger.With().
			Int(level, int(l)).
			Logger(),
		LogLevel: log.LogLevel,
	}
}

func (log *ZeroLoggerWrapper) Printf(format string, v ...interface{}) {
	log.logger.Log().Msg(fmt.Sprintf(format, v...))
}

func (log *ZeroLoggerWrapper) WithStr(k, v string) Logger {
	return &ZeroLoggerWrapper{
		logger:   log.logger.With().Str(k, v).Logger(),
		LogLevel: log.LogLevel,
	}
}

func (log *ZeroLoggerWrapper) WithInt(k string, v int) Logger {
	return &ZeroLoggerWrapper{
		logger:   log.logger.With().Int(k, v).Logger(),
		LogLevel: log.LogLevel,
	}
}

func (log *ZeroLoggerWrapper) WithFlo(k string, v float64) Logger {
	return &ZeroLoggerWrapper{
		logger:   log.logger.With().Float64(k, v).Logger(),
		LogLevel: log.LogLevel,
	}
}

func (log *ZeroLoggerWrapper) WithBoo(k string, v bool) Logger {
	return &ZeroLoggerWrapper{
		logger:   log.logger.With().Bool(k, v).Logger(),
		LogLevel: log.LogLevel,
	}
}

type NoopLogger struct {
}

func (*NoopLogger) Printf(string, ...interface{}) {

}

func (l *NoopLogger) Level(_ Level) Logger {
	return l
}

func (l *NoopLogger) WithStr(_, _ string) Logger {
	return l
}

func (l *NoopLogger) WithInt(_ string, _ int) Logger {
	return l
}

func (l *NoopLogger) WithFlo(_ string, _ float64) Logger {
	return l
}

func (l *NoopLogger) WithBoo(_ string, _ bool) Logger {
	return l
}
