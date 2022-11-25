package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	poseidon_server "github.com/zeus-fyi/olympus/poseidon/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := poseidon_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
