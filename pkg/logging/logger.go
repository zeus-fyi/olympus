package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

func SetLoggerLevel(level string) string {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	levelInt := strings.IntStringParser(level)
	loggingLevel := zerolog.Level(levelInt)
	zerolog.SetGlobalLevel(loggingLevel)

	log.Printf("logging is set to %s", loggingLevel)
	return loggingLevel.String()
}
