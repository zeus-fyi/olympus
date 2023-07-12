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
	MainnetRelay = "https://relay.flashbots.net"
	GoerliRelay  = "https://relay-goerli.flashbots.net"
	SepoliaRelay = "https://relay-sepolia.flashbots.net"

	BlocksAPI = "https://blocks.flashbots.net/v1/blocks"
	TxsAPI    = "https://blocks.flashbots.net/v1/transactions"

	FlashbotXHeader = "X-Flashbots-Signature"
)

type FlashbotsClient struct {
	resty_base.Resty
	W *web3_actions.Web3Actions
	flashbotsrpc.EthereumAPI
	*flashbotsrpc.FlashbotsRPC
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
