package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	artemis_server "github.com/zeus-fyi/olympus/artemis/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := artemis_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
