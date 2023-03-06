package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hydra_server "github.com/zeus-fyi/olympus/hydra/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if err := hydra_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
