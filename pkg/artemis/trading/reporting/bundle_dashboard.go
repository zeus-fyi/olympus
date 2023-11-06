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
	EventID        int                   `json:"eventID"`
	SubmissionTime string                `json:"submissionTime"`
	BundleHash     string                `json:"bundleHash"`
	TraderInfo     map[string]TraderInfo `json:"traderInfo"`
	BundleTxs      []Bundle              `json:"bundledTxs"`
}

type TraderInfo struct {
	TotalTxFees float64 `json:"totalTxFees"`
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
		sort.Slice(v, func(i, j int) bool {
			return v[i].TransactionIndex < v[j].TransactionIndex
		})
		if ds.Bundles[i].TraderInfo == nil {
			ds.Bundles[i].TraderInfo = make(map[string]TraderInfo)
		}
		for j, tx := range v {
			fees := (float64(tx.EffectiveGasPrice) * float64(tx.GasUsed)) / 1e18
			if _, ok := ds.Bundles[j].TraderInfo[tx.EthTx.From]; !ok {
				ds.Bundles[i].TraderInfo[tx.EthTx.From] = TraderInfo{
					TotalTxFees: fees,
				}
			} else {
				ds.Bundles[i].TraderInfo[tx.EthTx.From] = TraderInfo{
					TotalTxFees: ds.Bundles[i].TraderInfo[tx.EthTx.From].TotalTxFees + fees,
				}
			}
		}

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
