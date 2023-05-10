package artemis_flashbots

type FlashbotsBlocksV1Response struct {
	Blocks []struct {
		BlockNumber           int    `json:"block_number"`
		FeeRecipientEthDiff   string `json:"fee_recipient_eth_diff"`
		FeeRecipient          string `json:"fee_recipient"`
		EthSentToFeeRecipient string `json:"eth_sent_to_fee_recipient"`
		GasUsed               int    `json:"gas_used"`
		GasPrice              string `json:"gas_price"`
		Transactions          []struct {
			TransactionHash       string `json:"transaction_hash"`
			BundleType            string `json:"bundle_type"`
			TxIndex               int    `json:"tx_index"`
			BundleIndex           int    `json:"bundle_index"`
			BlockNumber           int    `json:"block_number"`
			EoaAddress            string `json:"eoa_address"`
			ToAddress             string `json:"to_address"`
			GasUsed               int    `json:"gas_used"`
			GasPrice              string `json:"gas_price"`
			EthSentToFeeRecipient string `json:"eth_sent_to_fee_recipient"`
			FeeRecipientEthDiff   string `json:"fee_recipient_eth_diff"`
		} `json:"transactions"`
	} `json:"blocks"`
	LatestBlockNumber int `json:"latest_block_number"`
}
