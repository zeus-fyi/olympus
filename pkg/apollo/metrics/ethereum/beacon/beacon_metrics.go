package apollo_beacon_prom_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	client_consts "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/constants"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"time"

	"github.com/zeus-fyi/olympus/pkg/iris/resty_base"
)

const (
	beaconConsensusSyncEndpoint = "/eth/v1/node/syncing"
	beaconExecSyncPayload       = `{"method":"eth_syncing","params":[],"id":1,"jsonrpc":"2.0"}`
)

// TODO embed this in the PrometheusMetrics struct
var (
//	BeaconConsensusSyncDistance = prometheus.NewGauge(prometheus.GaugeOpts{
//		Name: "ethereum_beacon_consensus_sync_distance",
//		Help: "How far behind head of chain is the beacon consensus client?",
//	})
//
//	BeaconConsensusSyncHeadSlot = prometheus.NewGauge(prometheus.GaugeOpts{
//		Name: "ethereum_beacon_consensus_sync_head_slot",
//		Help: "What slot is the beacon consensus client synced to?",
//	})
)

type ConsensusClientMetrics struct {
	ConsensusClientRestClient resty_base.Resty
	BeaconConsensusSyncStatus prometheus.Gauge
}

func NewConsensusClientMetrics(w apollo_metrics_workload_info.WorkloadInfo, bc BeaconConfig) ConsensusClientMetrics {
	m := ConsensusClientMetrics{
		ConsensusClientRestClient: resty_base.GetBaseRestyClient(bc.ConsensusClientSVC, ""),
	}

	m.BeaconConsensusSyncStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ethereum_beacon_consensus_sync_status_is_syncing",
		Help:        "Is the beacon consensus client syncing, or is it synced? 0 = syncing, 1 = synced",
		ConstLabels: nil,
	})
	return m
}

type ExecClientMetrics struct {
	ExecClientRestClient resty_base.Resty
	BeaconExecSyncStatus prometheus.Gauge
}

func NewExecClientMetrics(w apollo_metrics_workload_info.WorkloadInfo, bc BeaconConfig) ExecClientMetrics {
	m := ExecClientMetrics{
		ExecClientRestClient: resty_base.GetBaseRestyClient(bc.ExecClientSVC, ""),
	}
	m.BeaconExecSyncStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ethereum_beacon_exec_sync_status_is_syncing",
		Help:        "Is the beacon exec client syncing? 0 = syncing, 1 = synced",
		ConstLabels: nil,
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
		bm.BeaconConsensusSyncStatus,
		bm.BeaconExecSyncStatus,
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
		bm.BeaconConsensusClientSyncStatus()
		bm.BeaconExecClientSyncStatus()
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
}

func (bm *BeaconMetrics) BeaconExecClientSyncStatus() {
	log.Info().Msg("BeaconExecClientSyncStatus: getting sync status")
	bm.BeaconExecSyncStatus.Set(0)
	ss := client_consts.ExecClientSyncStatus{}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	resp, err := bm.ExecClientRestClient.R().
		SetResult(&ss).
		SetBody(beaconExecSyncPayload).Post("/")
	if err != nil {
		log.Err(err).Msgf("resp: %s", resp)
		return
	}
	// Should always return false if not syncing.
	if ss.Result == false {
		bm.BeaconConsensusSyncStatus.Set(1)
	}
}
