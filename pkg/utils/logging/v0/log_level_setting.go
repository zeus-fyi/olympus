package v0

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (l *LibV0) SetLoggerLevel(level zerolog.Level) string {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(level)
	log.Printf("logging is set to %s", level)
	return level.String()
}
