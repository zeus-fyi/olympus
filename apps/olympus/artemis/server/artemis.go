package artemis_server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	artemis_api_router "github.com/zeus-fyi/olympus/artemis/api"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	artemis_mev_tx_fetcher "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/mev"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
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
	log.Info().Msg("Artemis: Starting ArtemisMevWorkerGoerli")
	artemis_mev_tx_fetcher.ArtemisMevWorkerGoerli.Worker.RegisterWorker(c)
	err = artemis_mev_tx_fetcher.ArtemisMevWorkerGoerli.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisMevWorkerGoerli Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisMevWorkerGoerli Started")
	// mainnet
	log.Info().Msg("Artemis: Starting ArtemisEthereumMainnetTxBroadcastWorker")
	artemis_ethereum_transcations.ArtemisEthereumMainnetTxBroadcastWorker.Worker.RegisterWorker(c)
	err = artemis_ethereum_transcations.ArtemisEthereumMainnetTxBroadcastWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisEthereumMainnetTxBroadcastWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisEthereumMainnetTxBroadcastWorker Started")
	log.Info().Msg("Artemis: Starting ArtemisMevWorkerMainnet")
	artemis_mev_tx_fetcher.ArtemisMevWorkerMainnet.Worker.RegisterWorker(c)
	err = artemis_mev_tx_fetcher.ArtemisMevWorkerMainnet.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisMevWorkerMainnet Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisMevWorkerMainnet Started")

	log.Info().Msg("Artemis: Starting ArtemisActiveMevWorkerMainnet")
	artemis_mev_tx_fetcher.ArtemisActiveMevWorkerMainnet.Worker.RegisterWorker(c)
	err = artemis_mev_tx_fetcher.ArtemisActiveMevWorkerMainnet.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisActiveMevWorkerMainnet Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisActiveMevWorkerMainnet Started")

	log.Info().Msg("Artemis: Starting ArtemisMevWorkerMainnetHistoricalTxs")
	artemis_mev_tx_fetcher.ArtemisMevWorkerMainnetHistoricalTxs.Worker.RegisterWorker(c)
	err = artemis_mev_tx_fetcher.ArtemisMevWorkerMainnetHistoricalTxs.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Artemis: %s ArtemisMevWorkerMainnetHistoricalTxs Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Artemis: ArtemisMevWorkerMainnetHistoricalTxs Started")
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
	if env == "local" || env == "production-local" {
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowCredentials: true,
		}))
	} else {
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"https://cloud.zeus.fyi", "https://api.zeus.fyi", "https://hestia.zeus.fyi"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Headers", "X-CSRF-Token", "Accept-Encoding", "CloudCtxNsID"},
			AllowCredentials: true,
		}))
	}

	artemis_mev_tx_fetcher.InitUniswap(ctx, artemis_orchestration_auth.Bearer)
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
	Use:   "Orchestration for web3 txs, actions, and infra.",
	Short: "A microservice for web3 txs, actions, and infra orchestrations",
	Run: func(cmd *cobra.Command, args []string) {
		Artemis()
	},
}
