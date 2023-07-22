package artemis_flashbots

import (
	"context"
	"crypto/ecdsa"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

const (
	MainnetRelay      = "https://relay.flashbots.net"
	BlocknativeRelay  = "https://api.blocknative.com/v1/auction"
	BeaverRelay       = "https://rpc.beaverbuild.org/"
	Builder69         = "https://builder0x69.io"
	RsyncBuilder      = "https://rsync-builder.xyz/"
	TitanBuilder      = "https://rpc.titanbuilder.xyz"
	PayloadBuilder    = "https://rpc.payload.de"
	BuildAIBuilder    = "https://BuildAI.net"
	ManifolderBulder  = "https://api.securerpc.com/v1"
	NfactorialBuilder = "https://rpc.nfactorial.xyz/"
	EdenBuilder       = "https://api.edennetwork.io/v1/bundle"
	LighspeedBuilder  = "https://rpc.lightspeedbuilder.info/"
	EthBuilder        = "https://eth-builder.com"

	GoerliRelay  = "https://relay-goerli.flashbots.net"
	SepoliaRelay = "https://relay-sepolia.flashbots.net"

	BlocksAPI       = "https://blocks.flashbots.net/v1/blocks"
	TxsAPI          = "https://blocks.flashbots.net/v1/transactions"
	FlashbotXHeader = "X-Flashbots-Signature"
)

var Builders = []string{
	Builder69,
	RsyncBuilder,
	BlocknativeRelay,
	BeaverRelay,
	TitanBuilder,
	PayloadBuilder,
	BuildAIBuilder,
	NfactorialBuilder,
	EdenBuilder,
	LighspeedBuilder,
	EthBuilder,
	ManifolderBulder,
}

type FlashbotsClient struct {
	resty_base.Resty
	W *web3_actions.Web3Actions
	flashbotsrpc.EthereumAPI
	*flashbotsrpc.FlashbotsRPC
}

func InitFlashbotsClientForAdditionalBuilder(ctx context.Context, w *web3_actions.Web3Actions, builderRpc string) FlashbotsClient {
	rpc := flashbotsrpc.NewFlashbotsRPC(builderRpc)
	return FlashbotsClient{
		W:            w,
		FlashbotsRPC: rpc,
	}
}

func InitFlashbotsClient(ctx context.Context, w *web3_actions.Web3Actions) FlashbotsClient {
	rpc := flashbotsrpc.NewFlashbotsRPC(w.NodeURL)
	switch w.Network {
	case hestia_req_types.Mainnet:
		rpc = flashbotsrpc.New(MainnetRelay)
	case hestia_req_types.Goerli:
		rpc = flashbotsrpc.New(GoerliRelay)
	default:
		rpc = flashbotsrpc.New(MainnetRelay)
	}
	return FlashbotsClient{
		W:            w,
		FlashbotsRPC: rpc,
	}
}

func (f *FlashbotsClient) getPrivateKey() *ecdsa.PrivateKey {
	return f.W.Account.EcdsaPrivateKey()
}

func (f *FlashbotsClient) SendPrivateTx(ctx context.Context, privTx flashbotsrpc.FlashbotsSendPrivateTransactionRequest) (string, error) {
	resp, err := f.FlashbotsRPC.FlashbotsSendPrivateTransaction(f.getPrivateKey(), privTx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: SendPrivateTx")
		return resp, err
	}
	return resp, nil
}

func (f *FlashbotsClient) GetFlashbotsBlocksV1(ctx context.Context) (FlashbotsBlocksV1Response, error) {
	fbResp := FlashbotsBlocksV1Response{}
	_, err := f.R().SetResult(&fbResp).Get(BlocksAPI)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: GetFlashbotsBlocksV1")
		return fbResp, err
	}
	return fbResp, nil
}
