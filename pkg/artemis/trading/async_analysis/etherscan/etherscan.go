package artemis_etherscan

import "github.com/nanmu42/etherscan-api"

type Etherscan struct {
	*etherscan.Client
}

func NewMainnetEtherscanClient(apiKey string) Etherscan {
	client := etherscan.New(etherscan.Mainnet, apiKey)
	return Etherscan{
		Client: client,
	}
}
