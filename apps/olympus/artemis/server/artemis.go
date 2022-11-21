package artemis_server

import (
	"context"

	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	artemis_api_router "github.com/zeus-fyi/olympus/artemis/api"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

var cfg = Config{}

func Artemis() {
	cfg.Host = "0.0.0.0"
	srv := NewArtemisServer(cfg)
	// Echo instance
	srv.E = artemis_api_router.Routes(srv.E)
	ctx := context.Background()
	apps.Pg = apps.Db{}
	apps.Pg.InitPG(ctx, cfg.PGConnStr)

	// Middleware
	srv.E.Use(middleware.Logger())
	srv.E.Use(middleware.Recover())
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9004", "server port")
	Cmd.Flags().StringVar(&cfg.PGConnStr, "postgres-conn-str", "", "postgres connection string")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Orchestration for web3 & cloud",
	Short: "A microservice for orchestrations",
	Run: func(cmd *cobra.Command, args []string) {
		Artemis()
	},
}
