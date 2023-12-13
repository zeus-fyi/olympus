package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hephaestus_server "github.com/zeus-fyi/olympus/hephaestus/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := hephaestus_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
