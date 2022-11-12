package web3_client

import (
	"github.com/gochain/web3"
)

type Web3Client struct {
	NodeURL string
	Network string
	web3.Client
}

func NewClient(nodeURL string) Web3Client {
	return Web3Client{NodeURL: nodeURL}
}

func (w *Web3Client) Dial() {
	r, err := web3.Dial(w.NodeURL)
	if err != nil {
		panic(err)
	}
	w.Client = r
}
