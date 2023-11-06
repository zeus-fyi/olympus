package artemis_mev_transcations

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var ArtemisMevClientMainnet web3_client.Web3Client

func InitArtemisMevClientMainnet(ctx context.Context) {
	log.Info().Msg("Artemis: InitArtemisMevClientMainnet")
	cfg := artemis_network_cfgs.ArtemisEthereumMainnet
	if len(cfg.NodeURL) == 0 || cfg.Account == nil {
		err := errors.New("missing configs")
		log.Ctx(ctx).Panic().Err(err).Interface("nodeUrl", cfg.NodeURL).Interface("account", cfg.Account.PublicKey()).Msg("InitArtemisMevClientMainnet failed")
		misc.DelayedPanic(err)
	}
	defaultProxyUrl := "https://iris.zeus.fyi/v1/router"
	ArtemisMevClientMainnet = web3_client.NewWeb3Client(defaultProxyUrl, cfg.Account)
	if len(artemis_orchestration_auth.Bearer) == 0 {
		panic("missing bearer token")
	}
	ArtemisMevClientMainnet.AddBearerToken(artemis_orchestration_auth.Bearer)
	ArtemisMevClientMainnet.AddDefaultEthereumMainnetTableHeader()
	log.Info().Msg("Artemis: InitArtemisMevClientMainnet Succeeded")
}

var ArtemisMevClientGoerli web3_client.Web3Client

func InitArtemisMevClientGoerli(ctx context.Context) {
	log.Info().Msg("Artemis: InitArtemisMevClientGoerli")
	cfg := artemis_network_cfgs.ArtemisEthereumGoerli
	if len(cfg.NodeURL) == 0 || cfg.Account == nil {
		err := errors.New("missing configs")
		log.Ctx(ctx).Panic().Err(err).Interface("nodeUrl", cfg.NodeURL).Interface("account", cfg.Account.PublicKey()).Msg("InitArtemisMevClientGoerli failed")
		misc.DelayedPanic(err)
	}
	ArtemisMevClientGoerli = web3_client.NewWeb3Client(cfg.NodeURL, cfg.Account)
	log.Info().Msg("Artemis: InitArtemisEthereumGoerliClient ArtemisMevClientGoerli")
}

func InitWeb3Clients(ctx context.Context) {
	log.Ctx(ctx).Info().Msg("Artemis: Init MEV Web3Clients")
	InitArtemisMevClientMainnet(ctx)
	InitArtemisMevClientGoerli(ctx)
}
