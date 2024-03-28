package logger

import (
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/config"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	Level zerolog.Level
}

func New(config config.LoggerConfig) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	level := getLevel(config.Level)
	zerolog.SetGlobalLevel(level)
	return &Logger{
		Level: level,
	}
}

func getLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l Logger) Fatal(msg string) {
	log.Fatal().Msg(msg)
}

func (l Logger) Error(msg string) {
	log.Error().Msg(msg)
}

func (l Logger) Warn(msg string) {
	log.Warn().Msg(msg)
}

func (l Logger) Info(msg string) {
	log.Info().Msg(msg)
}

func (l Logger) Debug(msg string) {
	log.Debug().Msg(msg)
}
