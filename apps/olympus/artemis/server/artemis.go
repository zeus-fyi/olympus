package artemis_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	artemis_api_router "github.com/zeus-fyi/olympus/artemis/api"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var (
	cfg             = Config{}
	temporalAuthCfg temporal_auth.TemporalAuth
	env             string
	authKeysCfg     auth_keys_config.AuthKeysCfg
)

func Artemis() {
	cfg.Host = "0.0.0.0"
	srv := NewArtemisServer(cfg)
	// Echo instance
	ctx := context.Background()
	SetConfigByEnv(ctx, env)

	// goerli
	log.Info().Msg("Artemis: Starting ArtemisEthereumGoerliTxBroadcastWorker")
	c := artemis_ethereum_transcations.ArtemisEthereumGoerliTxBroadcastWorker.ConnectTemporalClient()
	defer c.Close()
	artemis_ethereum_transcations.ArtemisEthereumGoerliTxBroadcastWorker.Worker.RegisterWorker(c)
	err := artemis_ethereum_transcations.ArtemisEthereumGoerliTxBroadcastWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisEthereumGoerliTxBroadcastWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisEthereumGoerliTxBroadcastWorker Started")
	// mainnet
	log.Info().Msg("Artemis: Starting ArtemisEthereumMainnetTxBroadcastWorker")
	artemis_ethereum_transcations.ArtemisEthereumMainnetTxBroadcastWorker.Worker.RegisterWorker(c)
	err = artemis_ethereum_transcations.ArtemisEthereumMainnetTxBroadcastWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisEthereumMainnetTxBroadcastWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisEthereumMainnetTxBroadcastWorker Started")
	// ephemeral
	log.Info().Msg("Artemis: Starting ArtemisEthereumEphemeralTxBroadcastWorker")
	artemis_ethereum_transcations.ArtemisEthereumEphemeralTxBroadcastWorker.Worker.RegisterWorker(c)
	err = artemis_ethereum_transcations.ArtemisEthereumEphemeralTxBroadcastWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisEthereumEphemeralTxBroadcastWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisEthereumEphemeralTxBroadcastWorker Started")

	// Start server
	log.Info().Msg("Artemis: Starting Server")
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
