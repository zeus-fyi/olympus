package artemis_network_cfgs

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/configs"
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
	Ephemery = "ephemery"
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
	ArtemisEthereumMainnet         = NewArtemisConfig(Ethereum, Mainnet)
	ArtemisEthereumMainnetQuiknode = NewArtemisConfig(Ethereum, Mainnet)
	ArtemisEthereumGoerli          = NewArtemisConfig(Ethereum, Goerli)
	ArtemisEthereumEphemeral       = NewArtemisConfig(Ethereum, Ephemery)
	GlobalArtemisConfigs           = []*ArtemisConfig{&ArtemisEthereumMainnet, &ArtemisEthereumGoerli, &ArtemisEthereumEphemeral}
)

func (b *BeaconNetwork) GetBeaconSecretKey() string {
	return "secrets/" + strings.Join([]string{b.Service, b.Protocol, b.Network, "beacon"}, ".") + ".txt"
}

func (b *BeaconNetwork) GetBeaconWalletKey() string {
	return "secrets/" + strings.Join([]string{b.Service, b.Protocol, b.Network, "ecdsa"}, ".") + ".txt"
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

	ArtemisEthereumEphemeral.NodeURL = tc.EphemeralNodeUrl
	ArtemisEthereumGoerli.AddAccountFromHexPk(ctx, tc.ArtemisEphemeralEcdsaKey)
}
