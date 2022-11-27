package artemis_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_api_router "github.com/zeus-fyi/olympus/artemis/api"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var cfg = Config{}
var nodeURL string
var temporalAuthCfg temporal_auth.TemporalAuth
var env string
var artemisKey string
var authKeysCfg auth_keys_config.AuthKeysCfg

func Artemis() {
	cfg.Host = "0.0.0.0"
	srv := NewArtemisServer(cfg)
	// Echo instance
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
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		artemisKey = sw.ArtemisEcdsaKeys.Goerli
		cfg.PGConnStr = sw.PostgresAuth
		nodeURL = sw.GoerliNodeUrl
	case "production-local":
		log.Info().Msg("Artemis: production local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuth
		nodeURL = tc.GoerliNodeUrl
		artemisKey = tc.ArtemisGoerliEcdsaKey
	case "local":
		log.Info().Msg("Artemis: local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuth
		nodeURL = tc.GoerliNodeUrl
		artemisKey = tc.ArtemisGoerliEcdsaKey
	}
	log.Info().Msg("Artemis: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)

	log.Info().Msg("Artemis: connection client to web3 client")
	artemis, err := accounts.ParsePrivateKey(artemisKey)
	if err != nil {
		log.Info().Msg("Artemis: ethereum account failed to load")
		misc.DelayedPanic(err)
	}
	artemis_ethereum_transcations.InitArtemisEthereumClient(nodeURL, artemis)
	log.Info().Msgf("Artemis %s temporal auth and init procedure starting", env)
	artemis_ethereum_transcations.InitTxBroadcastWorker(temporalAuthCfg)

	// Start server
	srv.E = artemis_api_router.Routes(srv.E)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9004", "server port")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")
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
