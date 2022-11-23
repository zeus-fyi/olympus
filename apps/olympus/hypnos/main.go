package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hypnos_server "github.com/zeus-fyi/olympus/hypnos/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := hypnos_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
