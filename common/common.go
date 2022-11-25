package common

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerWrapper struct {
	log *zerolog.Logger
}

func Logger(ctx context.Context) *LoggerWrapper {
	return &LoggerWrapper{log: log.Ctx(ctx)}
}

func (l *LoggerWrapper) Infof(format string, v ...interface{}) {
	l.log.Info().Msgf(format, v)
}

func (l *LoggerWrapper) Debugf(format string, v ...interface{}) {
	l.log.Debug().Msgf(format, v)
}

func (l *LoggerWrapper) Errorf(format string, v ...interface{}) {
	l.log.Error().Msgf(format, v)
}

func (l *LoggerWrapper) Write(p []byte) (n int, err error) {
	l.log.Debug().Msgf(string(p))
	return len(p), nil
}
