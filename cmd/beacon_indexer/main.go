package main

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/server"
)

func main() {
	if err := server.ApiCmd.Execute(); err != nil {
		log.Err(err)
	}

}
