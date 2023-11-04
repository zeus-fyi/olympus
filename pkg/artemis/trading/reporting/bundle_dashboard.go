package artemis_reporting

/*
type BundlesGroup struct {
	Map map[string][]Bundle
}

type Bundle struct {
	artemis_autogen_bases.EthMevBundleProfit
	artemis_autogen_bases.EthTx
	artemis_autogen_bases.EthTxGas
	artemis_autogen_bases.EthTxReceipts
	artemis_autogen_bases.EthMempoolMevTx
	TradeExecutionFlow *web3_client.TradeExecutionFlow
}
*/

type BundleSummary struct {
	EventID    int      `json:"eventID"`
	BundleHash string   `json:"bundleHash"`
	BundleTxs  []Bundle `json:"bundledTxs"`
}

type BundleDashboardInfo struct {
	Bundles []BundleSummary `json:"bundles"`
}

func (b *BundlesGroup) GetDashboardInfo() BundleDashboardInfo {
	ds := BundleDashboardInfo{
		Bundles: make([]BundleSummary, len(b.bundleHashOrder)),
	}
	for ind, hash := range b.bundleHashOrder {
		ds.Bundles[ind].EventID = b.bundleHashToId[hash]
		ds.Bundles[ind].BundleHash = hash
		ds.Bundles[ind].BundleTxs = b.Map[hash]
	}
	return ds
}
