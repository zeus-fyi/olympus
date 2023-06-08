package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

const (
	QuikNodeSecret = "secrets/artemis.ethereum.mainnet.quiknode.txt"
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

	artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.NodeURL = secrets.MustReadSecret(ctx, inMemSecrets, QuikNodeSecret)
	artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.Account = artemis_network_cfgs.ArtemisEthereumMainnet.Account
	log.Info().Msg("Artemis: InitArtemisEthereum done")
	return
}
