package main

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/autok8s/server"
)

func main() {
	if err := server.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
