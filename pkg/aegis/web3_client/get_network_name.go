package web3_client

import "context"

const (
	Mainnet            = "Mainnet"
	Goerli             = "Goerli"
	MainnetGenesisHash = "0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"
	GoerliGenesisHash  = "0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a"
)

func (w *Web3Client) GetNetworkName(ctx context.Context) (string, error) {
	w.Dial()
	defer w.Close()
	id, err := w.GetID(ctx)
	if err != nil {
		return "", err
	}
	switch id.GenesisHash.Hex() {
	case GoerliGenesisHash:
		w.Network = Goerli
	case MainnetGenesisHash:
		w.Network = Mainnet
	}
	return w.Network, nil
}
