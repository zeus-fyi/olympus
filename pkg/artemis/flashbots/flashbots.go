package artemis_flashbots

import (
	"context"
	"crypto/ecdsa"

	"github.com/gochain/gochain/v4/common/hexutil"
	"github.com/gochain/gochain/v4/crypto"
	"github.com/metachris/flashbotsrpc"
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
}

func InitFlashbotsClient(ctx context.Context, nodeUrl, network string, acc *accounts.Account) FlashbotsClient {
	w := web3_actions.NewWeb3ActionsClientWithAccount(nodeUrl, acc)
	client := FlashbotsClient{resty_base.GetBaseRestyClient("", ""), w, flashbotsrpc.New(nodeUrl)}
	w.Network = network
	switch network {
	case hestia_req_types.Mainnet:
		client.SetBaseURL(MainnetRelay)
	case hestia_req_types.Goerli:
		client.SetBaseURL(GoerliRelay)
	}
	return client
}

func (f *FlashbotsClient) CreateFlashbotsHeader(signature []byte, privateKey *ecdsa.PrivateKey) string {
	return crypto.PubkeyToAddress(privateKey.PublicKey).Hex() +
		":" + hexutil.Encode(signature)
}
