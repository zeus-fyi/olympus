package athena_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	athena_router "github.com/zeus-fyi/olympus/athena/api"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

var cfg = Config{}
var authKeysCfg auth_keys_config.AuthKeysCfg
var env string
var dataDir structs.Path

func Athena() {
	ctx := context.Background()
	cfg.Host = "0.0.0.0"
	srv := NewAthenaServer(cfg)

	log.Info().Msgf("Athena: %s auth procedure starting", env)

	switch env {
	case "production":
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		dataDir.DirOut = "../"
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.LocalDbPgconn
		dataDir.DirOut = "../"
	}
	log.Info().Msg("Athena: DigitalOceanS3AuthClient starting")
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)

	log.Info().Msg("Athena: PG connection starting")
	apps.Pg = apps.Db{}
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	srv.E = athena_router.Routes(srv.E, dataDir)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9003", "server port")
	Cmd.Flags().StringVar(&cfg.PGConnStr, "postgres-conn-str", "", "postgres connection string")
	Cmd.Flags().StringVar(&dataDir.DirOut, "dataDir", "/data", "data directory location")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Web3 Middleware",
	Short: "A web3 infra middleware manager",
	Run: func(cmd *cobra.Command, args []string) {
		Athena()
	},
}
