package artemis_server

import (
	"context"

	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	artemis_api_router "github.com/zeus-fyi/olympus/artemis/api"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var cfg = Config{}
var nodeURL string
var temporalAuthCfg temporal_auth.TemporalAuth
var env string

func Artemis() {
	cfg.Host = "0.0.0.0"
	srv := NewArtemisServer(cfg)
	// Echo instance
	srv.E = artemis_api_router.Routes(srv.E)
	ctx := context.Background()
	switch env {
	case "production":
		log.Info().Msg("Artemis: production auth procedure starting")
		temporalAuthCfg = temporal_auth.TemporalAuth{
			ClientCertPath:   "/etc/ssl/certs/ca.pem",
			ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
			Namespace:        "production-zeus.ngb72",
			HostPort:         "production-zeus.ngb72.tmprl.cloud:7233",
		}
	case "production-local":
		log.Info().Msg("Artemis: production local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuth
	case "local":
		log.Info().Msg("Artemis: local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		temporalAuthCfg = tc.ProdLocalTemporalAuth
	}
	log.Info().Msg("Artemis: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)

	log.Info().Msg("Artemis: connection client to web3 client")
	artemis_ethereum_transcations.InitArtemisEthereumClient(nodeURL)

	log.Info().Msgf("Artemis %s temporal auth and init procedure starting", env)
	artemis_ethereum_transcations.InitTxBroadcastWorker(temporalAuthCfg)
	// TODO set bearer

	// Middleware
	srv.E.Use(middleware.Logger())
	srv.E.Use(middleware.Recover())
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9004", "server port")
	Cmd.Flags().StringVar(&cfg.PGConnStr, "postgres-conn-str", "postgresql://localhost/postgres?user=postgres&password=postgres", "postgres connection string")
	Cmd.Flags().StringVar(&nodeURL, "nodeURL", "", "node to use for client")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Orchestration for web3 interactions, cloud infra, and devops.",
	Short: "A microservice for orchestrations",
	Run: func(cmd *cobra.Command, args []string) {
		Artemis()
	},
}
