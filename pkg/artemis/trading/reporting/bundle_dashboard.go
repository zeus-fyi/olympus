package artemis_reporting

import (
	"sort"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

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
	EventID        int      `json:"eventID"`
	SubmissionTime string   `json:"submissionTime"`
	BundleHash     string   `json:"bundleHash"`
	BundleTxs      []Bundle `json:"bundledTxs"`
}

type BundleDashboardInfo struct {
	TopKTokens []string        `json:"topKTokens"`
	Bundles    []BundleSummary `json:"bundles"`
}

func (b *BundlesGroup) GetDashboardInfo() BundleDashboardInfo {
	ds := BundleDashboardInfo{
		Bundles: make([]BundleSummary, len(b.Map)),
	}

	ts := chronos.Chronos{}
	i := 0
	for hash, v := range b.Map {
		if len(v) < 2 {
			continue
		}
		ds.Bundles[i].EventID = b.MapHashToEventID[hash]
		ds.Bundles[i].SubmissionTime = ts.ConvertUnixTimeStampToDate(ds.Bundles[i].EventID).String()
		ds.Bundles[i].BundleHash = hash
		ds.Bundles[i].BundleTxs = v
		i++
	}
	// Truncate the slice to the number of elements actually set.
	if i < len(ds.Bundles) {
		ds.Bundles = ds.Bundles[:i]
	}
	sort.Slice(ds.Bundles, func(i, j int) bool {
		return ds.Bundles[i].EventID > ds.Bundles[j].EventID
	})
	return ds
}
