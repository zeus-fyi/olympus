package artemis_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	artemis_api_router "github.com/zeus-fyi/olympus/artemis/api"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var cfg = Config{}
var temporalAuthCfg temporal_auth.TemporalAuth
var env string
var authKeysCfg auth_keys_config.AuthKeysCfg

func Artemis() {
	cfg.Host = "0.0.0.0"
	srv := NewArtemisServer(cfg)
	// Echo instance
	ctx := context.Background()
	SetConfigByEnv(ctx, env)

	log.Info().Msg("Artemis: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Artemis: PG connection succeeded")

	// Start server
	srv.E = artemis_api_router.Routes(srv.E)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9004", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Orchestration for web3 txs.",
	Short: "A microservice for web3 tx orchestrations",
	Run: func(cmd *cobra.Command, args []string) {
		Artemis()
	},
}
