package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	olympus_downloader_init "github.com/zeus-fyi/olympus/downloader/downloader"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := olympus_downloader_init.Cmd.Execute(); err != nil {
		log.Err(err)
	}
}
