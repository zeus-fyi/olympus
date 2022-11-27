package artemis_network_cfgs

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

type ArtemisConfig struct {
	*accounts.Account
	BeaconNetwork
}

type ArtemisConfigs []ArtemisConfig

type BeaconNetwork struct {
	Service  string
	Protocol string
	Network  string
	NodeURL  string
}

const (
	Artemis  = "artemis"
	Mainnet  = "mainnet"
	Goerli   = "goerli"
	Ethereum = "ethereum"
)

func NewArtemisConfig(protocol, network string) ArtemisConfig {
	cfg := ArtemisConfig{
		BeaconNetwork: BeaconNetwork{
			Service:  Artemis,
			Protocol: protocol,
			Network:  network,
		},
	}
	return cfg
}

var (
	ArtemisEthereumMainnet = NewArtemisConfig(Ethereum, Mainnet)
	ArtemisEthereumGoerli  = NewArtemisConfig(Ethereum, Goerli)
	GlobalArtemisConfigs   = []ArtemisConfig{ArtemisEthereumMainnet, ArtemisEthereumGoerli}
)

func (b *BeaconNetwork) GetBeaconSecretKey() string {
	return "secrets/" + strings.Join([]string{b.Service, b.Protocol, b.Network}, ".") + "txt"
}

func (b *BeaconNetwork) GetBeaconWalletKey() string {
	return "secrets/" + strings.Join([]string{b.Service, b.Protocol, b.Network, "ecdsa"}, ".") + "txt"
}

func InitArtemisEthereum(ctx context.Context, inMemSecrets memfs.MemFS, secrets auth_startup.SecretsWrapper) {
	for _, cfg := range GlobalArtemisConfigs {
		cfg.NodeURL = secrets.ReadSecret(ctx, inMemSecrets, cfg.GetBeaconSecretKey())
		key := secrets.ReadSecret(ctx, inMemSecrets, cfg.GetBeaconWalletKey())
		cfg.AddAccountFromHexPk(ctx, key)
	}
	return
}

func (a *ArtemisConfig) AddAccountFromHexPk(ctx context.Context, key string) {
	acc, err := accounts.ParsePrivateKey(key)
	if err != nil {
		log.Ctx(ctx).Panic().Interface("AddAccountFromHexPk", a.BeaconNetwork).Msg("ArtemisConfig: ParsePrivateKey")
		panic(err)
	}
	a.Account = acc
}

func InitArtemisLocalTestConfigs() {
	tc := configs.InitLocalTestConfigs()
	ArtemisEthereumMainnet.NodeURL = tc.MainnetNodeUrl
	ctx := context.Background()
	ArtemisEthereumMainnet.AddAccountFromHexPk(ctx, tc.ArtemisMainnetEcdsaKey)

	ArtemisEthereumGoerli.NodeURL = tc.GoerliNodeUrl
	ArtemisEthereumGoerli.AddAccountFromHexPk(ctx, tc.ArtemisGoerliEcdsaKey)
}
