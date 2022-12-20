package artemis_ethereum_transcations

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var ArtemisEthereumBroadcastTxClient web3_client.Web3Client

func InitArtemisEthereumClient(ctx context.Context) {
	log.Info().Msg("Artemis: InitArtemisEthereumClient")
	cfg := artemis_network_cfgs.ArtemisEthereumMainnet
	if len(cfg.NodeURL) == 0 || cfg.Account == nil {
		err := errors.New("missing configs")
		log.Ctx(ctx).Panic().Err(err).Interface("nodeUrl", cfg.NodeURL).Interface("account", cfg.Account.PublicKey()).Msg("InitArtemisEthereumClient failed")
		misc.DelayedPanic(err)
	}
	ArtemisEthereumBroadcastTxClient = web3_client.NewWeb3Client(cfg.NodeURL, cfg.Account)
	log.Info().Msg("Artemis: InitArtemisEthereumClient Succeeded")
}

var ArtemisEthereumGoerliBroadcastTxClient web3_client.Web3Client

func InitArtemisEthereumGoerliClient(ctx context.Context) {
	log.Info().Msg("Artemis: InitArtemisEthereumGoerliClient")
	cfg := artemis_network_cfgs.ArtemisEthereumGoerli
	if len(cfg.NodeURL) == 0 || cfg.Account == nil {
		err := errors.New("missing configs")
		log.Ctx(ctx).Panic().Err(err).Interface("nodeUrl", cfg.NodeURL).Interface("account", cfg.Account.PublicKey()).Msg("InitArtemisEthereumGoerliClient failed")
		misc.DelayedPanic(err)
	}
	ArtemisEthereumGoerliBroadcastTxClient = web3_client.NewWeb3Client(cfg.NodeURL, cfg.Account)
	log.Info().Msg("Artemis: InitArtemisEthereumGoerliClient Succeeded")
}

var ArtemisEthereumEphemeralBroadcastTxClient web3_client.Web3Client

func InitArtemisEthereumEphemeralClient(ctx context.Context) {
	log.Info().Msg("Artemis: InitArtemisEthereumEphemeralClient")
	cfg := artemis_network_cfgs.ArtemisEthereumEphemeral
	if len(cfg.NodeURL) == 0 || cfg.Account == nil {
		err := errors.New("missing configs")
		log.Ctx(ctx).Panic().Err(err).Interface("nodeUrl", cfg.NodeURL).Interface("account", cfg.Account.PublicKey()).Msg("InitArtemisEthereumEphemeralClient failed")
		misc.DelayedPanic(err)
	}
	ArtemisEthereumEphemeralBroadcastTxClient = web3_client.NewWeb3Client(cfg.NodeURL, cfg.Account)
	log.Info().Msg("Artemis: InitArtemisEthereumGoerliClient Succeeded")
}

func InitWeb3Clients(ctx context.Context) {
	log.Ctx(ctx).Info().Msg("Artemis: InitWeb3Clients")
	InitArtemisEthereumClient(ctx)
	InitArtemisEthereumGoerliClient(ctx)
	InitArtemisEthereumEphemeralClient(ctx)
}
