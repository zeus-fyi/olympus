package artemis_flashbots

import (
	"context"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	"github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

const (
	MainnetRelay = "https://relay.flashbots.net"
	GoerliRelay  = "https://relay-goerli.flashbots.net"
	SepoliaRelay = "https://relay-sepolia.flashbots.net"

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
	client := FlashbotsClient{Resty: resty_base.GetBaseRestyClient("", ""), Web3Actions: w, EthereumAPI: flashbotsrpc.New(nodeUrl), FlashbotsRPC: rpc}
	w.Network = network
	switch network {
	case hestia_req_types.Mainnet:
		client.SetBaseURL(MainnetRelay)
	case hestia_req_types.Goerli:
		client.SetBaseURL(GoerliRelay)
	}
	return client
}

func (f *FlashbotsClient) SendBundle(ctx context.Context, bundle flashbotsrpc.FlashbotsSendBundleRequest) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	resp, err := f.FlashbotsRPC.FlashbotsSendBundle(f.EcdsaPrivateKey(), bundle)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsSendBundle")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	return resp, nil
}
