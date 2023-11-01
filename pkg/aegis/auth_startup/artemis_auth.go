package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

const (
	QuikNodeStreamWsSecret            = "secrets/artemis.ethereum.mainnet.quiknode.stream.ws.txt"
	QuikNodeSecret                    = "secrets/artemis.ethereum.mainnet.quiknode.txt"
	QuikNodeSecretLive                = "secrets/artemis.ethereum.mainnet.quiknode.live.txt"
	QuikNodeSecretLiveTest            = "secrets/artemis.ethereum.mainnet.quiknode.live.test.txt"
	QuiknodeHistoricalPrimarySecret   = "secrets/artemis.ethereum.mainnet.quiknode.historical.primary.txt"
	QuiknodeHistoricalSecondarySecret = "secrets/artemis.ethereum.mainnet.quiknode.historical.secondary.txt"
	QuiknodeHistoricalTertiarySecret  = "secrets/artemis.ethereum.mainnet.quiknode.historical.tertiary.txt"
)

func InitArtemisEthereum(ctx context.Context, inMemSecrets memfs.MemFS, secrets SecretsWrapper) {
	log.Info().Msg("Artemis: InitArtemisEthereum starting")
	for _, cfg := range artemis_network_cfgs.GlobalArtemisConfigs {
		beaconKey := cfg.GetBeaconSecretKey()
		beaconUrl := secrets.MustReadSecret(ctx, inMemSecrets, beaconKey)
		cfg.NodeURL = beaconUrl
		key := secrets.MustReadSecret(ctx, inMemSecrets, cfg.GetBeaconWalletKey())
		cfg.AddAccountFromHexPk(ctx, key)
	}

	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLiveTest.NodeURL = secrets.MustReadSecret(ctx, inMemSecrets, QuikNodeSecretLiveTest)

	artemis_test_cache.InitLiveTestNetwork("https://iris.zeus.fyi/v1/router")

	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalPrimary.NodeURL = secrets.MustReadSecret(ctx, inMemSecrets, QuiknodeHistoricalPrimarySecret)
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalPrimary.Account = artemis_network_cfgs.ArtemisEthereumMainnet.Account

	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalSecondary.NodeURL = secrets.MustReadSecret(ctx, inMemSecrets, QuiknodeHistoricalSecondarySecret)
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalSecondary.Account = artemis_network_cfgs.ArtemisEthereumMainnet.Account

	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalTertiary.NodeURL = secrets.MustReadSecret(ctx, inMemSecrets, QuiknodeHistoricalTertiarySecret)
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalTertiary.Account = artemis_network_cfgs.ArtemisEthereumMainnet.Account

	artemis_network_cfgs.ArtemisQuicknodeStreamWebsocket = secrets.MustReadSecret(ctx, inMemSecrets, QuikNodeStreamWsSecret)
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL = secrets.MustReadSecret(ctx, inMemSecrets, QuikNodeSecretLive)
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.Account = artemis_network_cfgs.ArtemisEthereumMainnet.Account

	artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.NodeURL = secrets.MustReadSecret(ctx, inMemSecrets, QuikNodeSecret)
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.Account = artemis_network_cfgs.ArtemisEthereumMainnet.Account
	log.Info().Msg("Artemis: InitArtemisEthereum done")
	return
}
