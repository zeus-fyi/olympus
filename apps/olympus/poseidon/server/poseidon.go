package poseidon_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	v1_poseidon "github.com/zeus-fyi/olympus/poseidon/api/v1"
	poseidon_pkg "github.com/zeus-fyi/olympus/poseidon/pkg"
)

var cfg = Config{}
var authKeysCfg auth_keys_config.AuthKeysCfg
var env string

func Poseidon() {
	cfg.Host = "0.0.0.0"
	ctx := context.Background()
	srv := NewPoseidonServer(cfg)
	// Echo instance
	srv.E = v1_poseidon.Routes(srv.E)

	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.LocalDbPgconn
	}

	log.Info().Msg("Poseidon: DigitalOceanS3AuthClient starting")
	poseidon_pkg.InitPoseidonReader(auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg))
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9006", "server port")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Downloading service for blockchain data",
	Short: "Chain download service",
	Run: func(cmd *cobra.Command, args []string) {
		Poseidon()
	},
}
