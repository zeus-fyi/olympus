package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	olympus_snapshot_init "github.com/zeus-fyi/olympus/downloader/startup_procedures"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := olympus_snapshot_init.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
