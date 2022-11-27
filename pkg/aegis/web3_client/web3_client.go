package web3_client

import (
	"github.com/zeus-fyi/gochain/web3/client"
	ecdsa_signer "github.com/zeus-fyi/olympus/pkg/aegis/ecdsa"
)

type Web3Client struct {
	NodeURL string
	Network string
	client.Client
	ecdsa_signer.EcdsaSigner
}

func NewClient(nodeURL string) Web3Client {
	return Web3Client{NodeURL: nodeURL}
}

func NewClientWithSigner(nodeURL string, ecdsaSigner ecdsa_signer.EcdsaSigner) Web3Client {
	wc := Web3Client{NodeURL: nodeURL, EcdsaSigner: ecdsaSigner}
	return wc
}

func (w *Web3Client) Dial() {
	r, err := client.Dial(w.NodeURL)
	if err != nil {
		panic(err)
	}
	w.Client = r
}
