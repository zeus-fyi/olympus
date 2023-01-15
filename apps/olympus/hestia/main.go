package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hestia_server "github.com/zeus-fyi/olympus/hestia/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := hestia_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
