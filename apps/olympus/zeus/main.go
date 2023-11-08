package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	zeus_server "github.com/zeus-fyi/olympus/zeus/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := zeus_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
