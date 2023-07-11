package artemis_flashbots

import (
	"context"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
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
	web3_actions.Web3Actions
	flashbotsrpc.EthereumAPI
	*flashbotsrpc.FlashbotsRPC
}

func InitFlashbotsClient(ctx context.Context, nodeUrl, network string, acc *accounts.Account) FlashbotsClient {
	w := web3_actions.NewWeb3ActionsClientWithAccount(nodeUrl, acc)
	rpc := flashbotsrpc.NewFlashbotsRPC(nodeUrl)
	w.Network = network
	switch network {
	case hestia_req_types.Mainnet:
		rpc = flashbotsrpc.New(MainnetRelay)
	case hestia_req_types.Goerli:
		rpc = flashbotsrpc.New(GoerliRelay)
	default:
		rpc = flashbotsrpc.New(MainnetRelay)
	}
	return FlashbotsClient{
		Web3Actions:  w,
		FlashbotsRPC: rpc,
	}
}

func (f *FlashbotsClient) SendPrivateTx(ctx context.Context, privTx flashbotsrpc.FlashbotsSendPrivateTransactionRequest) (string, error) {
	resp, err := f.FlashbotsRPC.FlashbotsSendPrivateTransaction(f.EcdsaPrivateKey(), privTx)
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
