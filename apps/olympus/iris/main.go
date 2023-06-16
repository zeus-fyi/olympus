package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	iris_server "github.com/zeus-fyi/olympus/iris/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := iris_server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
