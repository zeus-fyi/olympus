package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hydra_choreography "github.com/zeus-fyi/olympus/choreography/hydra/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := hydra_choreography.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
