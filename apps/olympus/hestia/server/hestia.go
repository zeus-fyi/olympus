package server

import (
	"context"

	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	v1 "github.com/zeus-fyi/olympus/hestia/api/v1"
)

var cfg = Config{}

func Hestia() {
	srv := NewHestiaServer(cfg)
	// Echo instance
	srv.E = v1.Routes(srv.E)
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
	Cmd.Flags().StringVar(&cfg.Port, "port", "9002", "server port")
	Cmd.Flags().StringVar(&cfg.PGConnStr, "postgres-conn-str", "", "postgres connection string")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Storing Internal Data",
	Short: "A microservice for internal configurations",
	Run: func(cmd *cobra.Command, args []string) {
		Hestia()
	},
}
