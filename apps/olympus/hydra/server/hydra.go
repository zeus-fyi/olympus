package hydra_server

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	v1_hydra "github.com/zeus-fyi/olympus/hydra/api/v1"
	hydra_eth2_web3signer "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer"
	ethereum_slashing_protection_watermarking "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/slashing_protection"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	eth_validator_signature_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var (
	cfg             = Config{}
	temporalAuthCfg temporal_auth.TemporalAuth
	authKeysCfg     auth_keys_config.AuthKeysCfg
	env             string
	Workload        WorkloadInfo
)

const (
	Mainnet  = "mainnet"
	Ephemery = "ephemery"
)

type WorkloadInfo struct {
	zeus_common_types.CloudCtxNs
	ProtocolNetworkID int // eg. mainnet
}

func Hydra() {
	ctx := context.Background()
	cfg.Host = "0.0.0.0"
	srv := NewHydraServer(cfg)
	log.Ctx(ctx).Info().Msg("Hydra: Initializing configs by environment type")
	SetConfigByEnv(ctx, env)
	srv.E = v1_hydra.Routes(srv.E)

	log.Ctx(ctx).Info().Msg("Hydra: Starting Async Service Route Polling")
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID: Workload.ProtocolNetworkID,
	}
	artemis_validator_signature_service_routing.InitAsyncServiceAuthRoutePolling(ctx, vsi, Workload.CloudCtxNs)
	log.Ctx(ctx).Info().Msg("Hydra: Async Service Route Polling Started")

	log.Ctx(ctx).Info().Msg("Hydra: Starting Temporal Worker")
	switch ethereum_slashing_protection_watermarking.Network {
	case Mainnet:
		eth_validator_signature_requests.InitArtemisEthereumValidatorSignatureRequestsMainnetWorker(ctx, temporalAuthCfg)
		c := eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsMainnetWorker.ConnectTemporalClient()
		defer c.Close()
		eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsMainnetWorker.RegisterWorker(c)
		err := eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsMainnetWorker.Worker.Start()
		if err != nil {
			log.Fatal().Err(err).Msgf("Hydra: %s ArtemisEthereumValidatorSignatureRequestsMainnetWorker.Worker.Start failed", env)
			misc.DelayedPanic(err)
		}
	case Ephemery:
		eth_validator_signature_requests.InitArtemisEthereumValidatorSignatureRequestsEphemeryWorker(ctx, temporalAuthCfg)
		c := eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.ConnectTemporalClient()
		defer c.Close()
		eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.RegisterWorker(c)
		err := eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.Worker.Start()
		if err != nil {
			log.Fatal().Err(err).Msgf("Hydra: %s ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.Worker.Start failed", env)
			misc.DelayedPanic(err)
		}
	default:
		panic("unknown network")
	}
	log.Ctx(ctx).Info().Interface("network", ethereum_slashing_protection_watermarking.Network).Msg("Hydra: Temporal Worker Started")

	log.Ctx(ctx).Info().Msg("Hydra: Starting async priority message queues")
	hydra_eth2_web3signer.InitAsyncMessageQueues(ctx)
	log.Ctx(ctx).Info().Msg("Hydra: Async priority message queues started")

	log.Ctx(ctx).Info().Msg("Hydra: Starting server")
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9000", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().IntVar(&Workload.ProtocolNetworkID, "protocol-network-id", hestia_req_types.EthereumEphemeryProtocolNetworkID, "identifier for protocol and network")

	Cmd.Flags().StringVar(&Workload.CloudCtxNs.CloudProvider, "cloud-provider", "", "cloud-provider")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Context, "ctx", "", "context")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Namespace, "ns", "", "namespace")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Region, "region", "", "region")

	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Web3signing proxy router",
	Short: "Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		Hydra()
	},
}
