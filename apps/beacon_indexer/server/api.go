package server

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "github.com/zeus-fyi/olympus/beacon-indexer/api/v1"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_fetcher"
	"github.com/zeus-fyi/olympus/pkg/databases/postgres"
)

var (
	BeaconEndpointURL string
	PGConnStr         string
)

func Api() {
	// Echo instance
	e := echo.New()
	e = v1.Routes(e)
	postgres.Pg = postgres.Db{}
	postgres.Pg.InitPG(context.Background(), PGConnStr)
	beacon_fetcher.InitFetcherService(BeaconEndpointURL)
	// Start server
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	err := e.Start(":9000")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	viper.AutomaticEnv()
	ApiCmd.Flags().StringVar(&PGConnStr, "postgres-conn-str", "", "postgres connection string")
	ApiCmd.Flags().StringVar(&BeaconEndpointURL, "beacon-endpoint", "", "beacon endpoint url")
}

// ApiCmd represents the base command when called without any subcommands
var ApiCmd = &cobra.Command{
	Use:   "eth2-indexer",
	Short: "A Reporting Focused Indexer for Ethereum 2.0/Consensus Layer",
	Run: func(cmd *cobra.Command, args []string) {
		Api()
	},
}
