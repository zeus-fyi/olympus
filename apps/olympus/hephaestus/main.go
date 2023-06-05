package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hephaestus_build_actions "github.com/zeus-fyi/olympus/pkg/hephaestus"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := hephaestus_build_actions.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
