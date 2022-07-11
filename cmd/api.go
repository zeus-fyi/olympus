package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/beacon_fetcher"
)

var (
	BeaconEndpointURL string
	PGConnStr         string
)

func Api() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", health)

	postgres.Pg = postgres.Db{}
	postgres.Pg.InitPG(context.Background(), PGConnStr)
	beacon_fetcher.InitFetcherService(BeaconEndpointURL)
	// Start server
	err := e.Start(":9000")
	if err != nil {
		log.Fatal(err)
	}
}

// Handler
func health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
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
