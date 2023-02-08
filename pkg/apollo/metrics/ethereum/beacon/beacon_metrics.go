package apollo_beacon_prom_metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	client_consts "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/constants"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"path"
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
	BeaconConsensusSyncStatus prometheus.Gauge
}

func NewConsensusClientMetrics(w apollo_metrics_workload_info.WorkloadInfo) ConsensusClientMetrics {
	m := ConsensusClientMetrics{}
	m.BeaconConsensusSyncStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "ethereum_beacon_consensus_sync_status_is_syncing",
		Help:        "Is the beacon consensus client syncing, or is it synced? 0 = syncing, 1 = synced",
		ConstLabels: nil,
	})
	return m
}

type ExecClientMetrics struct {
	BeaconExecSyncStatus prometheus.Gauge
}

func NewExecClientMetrics(w apollo_metrics_workload_info.WorkloadInfo) ExecClientMetrics {
	m := ExecClientMetrics{}
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
		ExecClientSVC:       fmt.Sprintf("http://zeus-exec-client:%s", client_consts.GetAnyClientApiPorts("geth")),
		ConsensusClientSVC:  fmt.Sprintf("http://zeus-consensus-client:%s", client_consts.GetAnyClientApiPorts("lighthouse")),
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
		ConsensusClientMetrics: NewConsensusClientMetrics(w),
		ExecClientMetrics:      NewExecClientMetrics(w),
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
	bm.BeaconConsensusSyncStatus.Set(0)
	ss := client_consts.ConsensusClientSyncStatus{}

	resp, err := bm.R.R().
		SetResult(&ss).
		Get(path.Join(bm.ConsensusClientSVC, beaconConsensusSyncEndpoint))
	if err != nil {
		log.Err(err).Msgf("resp: %s", resp)
		return
	}
	if ss.Data.IsSyncing == false {
		bm.BeaconConsensusSyncStatus.Set(1)
	}
}

func (bm *BeaconMetrics) BeaconExecClientSyncStatus() {
	bm.BeaconExecSyncStatus.Set(0)
	ss := client_consts.ExecClientSyncStatus{}

	resp, err := bm.R.R().
		SetResult(&ss).
		SetBody(beaconExecSyncPayload).
		Post(bm.ExecClientSVC)
	if err != nil {
		log.Err(err).Msgf("resp: %s", resp)
		return
	}
	// Should always return false if not syncing.
	if ss.Result == false {
		bm.BeaconConsensusSyncStatus.Set(1)
	}
}
