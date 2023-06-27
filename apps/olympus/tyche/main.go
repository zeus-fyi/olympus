package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tyche_server "github.com/zeus-fyi/olympus/tyche/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if err := tyche_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
