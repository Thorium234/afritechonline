package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New returns a configured zerolog instance.
func New(env string) zerolog.Logger {
	level := zerolog.InfoLevel
	if env == "development" {
		level = zerolog.DebugLevel
		zerolog.SetGlobalLevel(level)
		// Pretty console output in development
		return log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	}
	zerolog.SetGlobalLevel(level)
	return log.With().Timestamp().Logger()
}
