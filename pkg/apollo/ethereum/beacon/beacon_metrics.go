package apollo_beacon_prom_metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	"github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	client_consts "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/constants"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

const (
	beaconConsensusSyncEndpoint = "/eth/v1/node/syncing"
	beaconExecSyncPayload       = `{"method":"eth_syncing","params":[],"id":1,"jsonrpc":"2.0"}`
	beaconExecBlockHeight       = `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}`
)

type ExecClientSyncStatusBlockHeight struct {
	JsonRPC string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

type ConsensusClientMetrics struct {
	ConsensusClientRestClient         resty_base.Resty
	BeaconConsensusSyncStatus         prometheus.Gauge
	BeaconConsensusClientSyncDistance prometheus.Gauge
}

func NewConsensusClientMetrics(w apollo_metrics_workload_info.WorkloadInfo, bc BeaconConfig) ConsensusClientMetrics {
	m := ConsensusClientMetrics{
		ConsensusClientRestClient: resty_base.GetBaseRestyClient(bc.ConsensusClientSVC, ""),
	}
	m.BeaconConsensusClientSyncDistance = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ethereum_beacon_consensus_sync_distance",
		Help: "How many slots is the beacon consensus client behind head of chain?",
	})
	m.BeaconConsensusSyncStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ethereum_beacon_consensus_sync_status_is_syncing",
		Help: "Is the beacon consensus client syncing, or is it synced? 0 = syncing, 1 = synced",
	})
	return m
}

type ExecClientMetrics struct {
	ExecClientRestClient      resty_base.Resty
	BeaconExecSyncStatus      prometheus.Gauge
	BeaconExecSyncBlockHeight prometheus.Gauge
}

func NewExecClientMetrics(w apollo_metrics_workload_info.WorkloadInfo, bc BeaconConfig) ExecClientMetrics {
	m := ExecClientMetrics{
		ExecClientRestClient: resty_base.GetBaseRestyClient(bc.ExecClientSVC, ""),
	}
	m.BeaconExecSyncBlockHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ethereum_beacon_exec_sync_status_block_height",
		Help: "Returns the current block number the client is on.",
	})
	m.BeaconExecSyncStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ethereum_beacon_exec_sync_status_is_syncing",
		Help: "Is the beacon exec client syncing? 0 = syncing, 1 = synced",
	})
	return m
}

type BeaconConfig struct {
	Name                string
	Network             string
	ConsensusClientName string
	ExecClientName      string
	BeaconURL           string
	ExecClientSVC       string
	ConsensusClientSVC  string
}

func HydraConfig(network string) BeaconConfig {
	return BeaconConfig{
		Name:                "hydra",
		Network:             network,
		ConsensusClientName: "lighthouse",
		ExecClientName:      "geth",
		BeaconURL:           "",
		ExecClientSVC:       "http://zeus-exec-client:8545",
		ConsensusClientSVC:  "http://zeus-consensus-client:5052",
	}
}

type BeaconMetrics struct {
	R resty_base.Resty // metricsServicesPollingClient
	BeaconConfig
	zeus_common_types.CloudCtxNs
	ConsensusClientMetrics
	ExecClientMetrics
}

func (bm *BeaconMetrics) GetMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		bm.BeaconConsensusClientSyncDistance,
		bm.BeaconConsensusSyncStatus,
		bm.BeaconExecSyncStatus,
		bm.BeaconExecSyncBlockHeight,
	}
}

func NewBeaconMetrics(w apollo_metrics_workload_info.WorkloadInfo, bc BeaconConfig, bearer string) BeaconMetrics {
	return BeaconMetrics{
		R:                      resty_base.GetBaseRestyClient("", bearer),
		BeaconConfig:           bc,
		CloudCtxNs:             w.CloudCtxNs,
		ConsensusClientMetrics: NewConsensusClientMetrics(w, bc),
		ExecClientMetrics:      NewExecClientMetrics(w, bc),
	}
}

func (bm *BeaconMetrics) PollMetrics(pollTime time.Duration) {
	ticker := time.Tick(pollTime)
	for range ticker {
		go bm.BeaconConsensusClientSyncStatus()
		go bm.BeaconExecClientSyncStatus()
	}
}

func (bm *BeaconMetrics) BeaconConsensusClientSyncStatus() {
	log.Info().Msg("BeaconConsensusClientSyncStatus: getting sync status")
	bm.BeaconConsensusSyncStatus.Set(0)
	ss := client_consts.ConsensusClientSyncStatus{}
	resp, err := bm.ConsensusClientRestClient.R().
		SetResult(&ss).
		Get(beaconConsensusSyncEndpoint)
	if err != nil {
		log.Err(err).Msgf("resp: %s", resp)
		return
	}
	if ss.Data.IsSyncing == false {
		bm.BeaconConsensusSyncStatus.Set(1)
	}
	parsedNumber, err := strconv.ParseFloat(ss.Data.SyncDistance, 64)
	if err != nil {
		log.Err(err).Msgf("parsedNumber err: %s", err)
		return
	}
	bm.BeaconConsensusClientSyncDistance.Set(parsedNumber)
}

func (bm *BeaconMetrics) BeaconExecClientSyncStatus() {
	log.Info().Msg("BeaconExecClientSyncStatus: getting sync status")
	bm.BeaconExecSyncStatus.Set(0)
	ss := client_consts.ExecClientSyncStatus{}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	resp, err := bm.ExecClientRestClient.R().
		SetHeaders(headers).
		SetResult(&ss).
		SetBody(beaconExecSyncPayload).Post("/")
	if err != nil {
		log.Err(err).Msgf("resp: %s", resp)
		return
	}
	// Should always return false if not syncing.
	if ss.Result == false {
		bm.BeaconExecSyncStatus.Set(1)
	}
	bh := ExecClientSyncBlockHeight{}
	resp, err = bm.ExecClientRestClient.R().
		SetHeaders(headers).
		SetResult(&bh).
		SetBody(beaconExecBlockHeight).Post("/")
	if err != nil {
		log.Err(err).Msgf("resp: %s", resp)
		return
	}
	blockNum, err := strconv.ParseInt(strings_filter.Trim0xPrefix(bh.Result), 16, 64)
	if err != nil {
		log.Err(err).Msgf("can't decode hex string result %s", bh.Result)
		return
	}
	floatValueBlockNum := float64(blockNum)
	bm.BeaconExecSyncBlockHeight.Set(floatValueBlockNum)
}

type ExecClientSyncBlockHeight struct {
	JsonRPC string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}
