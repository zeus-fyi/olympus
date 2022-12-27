package hera_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	v1_hera "github.com/zeus-fyi/olympus/hera/api/v1"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var cfg = Config{}
var authKeysCfg auth_keys_config.AuthKeysCfg
var temporalAuthCfg temporal_auth.TemporalAuth
var env string

func Hera() {
	ctx := context.Background()

	cfg.Host = "0.0.0.0"
	srv := NewHeraServer(cfg)
	// Echo instance
	srv.E = v1_hera.Routes(srv.E)

	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunHeraDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		temporalAuthCfg = temporal_auth.TemporalAuth{
			ClientCertPath:   "/etc/ssl/certs/ca.pem",
			ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
			Namespace:        "production-hera.ngb72",
			HostPort:         "production-hera.ngb72.tmprl.cloud:7233",
		}
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
	}

	log.Info().Msg("Hera: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9008", "server port")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")

}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Generation of keys, networks, code, and more.",
	Short: "Generation of keys, networks, code..",
	Run: func(cmd *cobra.Command, args []string) {
		Hera()
	},
}
