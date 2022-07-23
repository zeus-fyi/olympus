package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// shared logging
func init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	loggingLevel := zerolog.DebugLevel
	zerolog.SetGlobalLevel(loggingLevel)

	log.Printf("logging is set to %s", loggingLevel)
}
