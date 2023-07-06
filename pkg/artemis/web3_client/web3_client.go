package web3_client

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

type Web3Client struct {
	web3_actions.Web3Actions
}

func NewWeb3ClientFakeSigner(nodeUrl string) Web3Client {
	acc, _ := accounts.CreateAccount()
	w := web3_actions.NewWeb3ActionsClientWithAccount(nodeUrl, acc)
	return Web3Client{w}
}

func NewWeb3ClientWithRelay(relayProxyUrl, nodeUrl string, acc *accounts.Account) Web3Client {
	w := web3_actions.NewWeb3ActionsClientWithRelayProxy(relayProxyUrl, nodeUrl, acc)
	return Web3Client{w}
}

func NewWeb3Client(nodeUrl string, acc *accounts.Account) Web3Client {
	w := web3_actions.NewWeb3ActionsClientWithAccount(nodeUrl, acc)
	return Web3Client{w}
}
