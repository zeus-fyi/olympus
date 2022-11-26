package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	aegis_server "github.com/zeus-fyi/olympus/aegis/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := aegis_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
