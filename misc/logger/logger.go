package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(serviceName string) zerolog.Logger {
	fields := map[string]any{
		"service": serviceName,
	}

	log := zerolog.New(os.Stdout).With().Fields(fields).Timestamp().Logger()

	return log
}
