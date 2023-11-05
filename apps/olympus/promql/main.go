package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	promql_server "github.com/zeus-fyi/olympus/promql/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := promql_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
