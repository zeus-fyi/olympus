package hestia_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	v1hestia "github.com/zeus-fyi/olympus/hestia/api/v1"
)

var (
	cfg = Config{}
	env string
)

func Hestia() {
	cfg.Host = "0.0.0.0"
	srv := NewHestiaServer(cfg)
	// Echo instance
	srv.E = v1hestia.Routes(srv.E)
	ctx := context.Background()
	switch env {
	case "production":
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
	}
	log.Info().Msg("Hestia: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9002", "server port")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Storing Internal Data",
	Short: "A microservice for internal configurations",
	Run: func(cmd *cobra.Command, args []string) {
		Hestia()
	},
}
