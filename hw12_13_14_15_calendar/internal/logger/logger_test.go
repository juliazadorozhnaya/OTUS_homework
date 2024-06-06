package logger

import (
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/config"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		New(&config.LoggerConfig{
			Level: "info",
		})
		require.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())

		New(&config.LoggerConfig{
			Level: "error",
		})
		require.Equal(t, zerolog.ErrorLevel, zerolog.GlobalLevel())

		New(&config.LoggerConfig{
			Level: "warn",
		})
		require.Equal(t, zerolog.WarnLevel, zerolog.GlobalLevel())
	})
}
