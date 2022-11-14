package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	athena_server "github.com/zeus-fyi/olympus/athena/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := athena_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
