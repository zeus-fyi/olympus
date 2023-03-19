package olympus_beacon_cookbooks

import (
	"fmt"

	olympus_hydra_choreography_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra/choreography"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1Core "k8s.io/api/core/v1"
)

const (
	athena                = "athena"
	consensusClient       = "zeus-consensus-client"
	execClient            = "zeus-exec-client"
	lighthouseDockerImage = "sigp/lighthouse:v3.5.1"
	gethDockerImage       = "ethereum/client-go:v1.11.2"
	initSnapshots         = "init-snapshots"
)

var (
	EphemeralBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumEphemeralAthenaBeacon",
		CloudCtxNs:       GetBeaconCloudCtxNs(hestia_req_types.Ephemery),
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  GetConsensusClientComponentBaseConfig(hestia_req_types.Ephemery),
			"execClients":       GetExecClientComponentBaseConfig(hestia_req_types.Ephemery),
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	GoerliBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumGoerliAthenaBeacon",
		CloudCtxNs:       GetBeaconCloudCtxNs(hestia_req_types.Goerli),
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  GetConsensusClientComponentBaseConfig(hestia_req_types.Goerli),
			"execClients":       GetExecClientComponentBaseConfig(hestia_req_types.Goerli),
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	MainnetBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumMainnetAthenaBeacon",
		CloudCtxNs:       GetBeaconCloudCtxNs(hestia_req_types.Mainnet),
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  GetConsensusClientComponentBaseConfig(hestia_req_types.Mainnet),
			"execClients":       GetExecClientComponentBaseConfig(hestia_req_types.Mainnet),
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	ConsensusClientComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"lighthouseAthena": ConsensusClientSkeletonBaseConfig,
		},
	}
	ConsensusClientSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ConsensusClientChartPath,
	}
	ConsensusClientChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/beacons/infra/consensus_client",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "lighthouseAthena", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
	ExecClientComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"gethAthena": ExecClientSkeletonBaseConfig,
		},
	}
	ExecClientSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ExecClientChartPath,
	}
	ExecClientChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/beacons/infra/exec_client",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "gethAthena", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)

func GetBeaconCloudCtxNs(network string) zeus_common_types.CloudCtxNs {
	beaconCloudCtxNs := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     fmt.Sprintf("athena-beacon-%s", network), // set with your own namespace
		Env:           "production",
	}
	return beaconCloudCtxNs
}

func GetConsensusClientComponentBaseConfig(network string) zeus_cluster_config_drivers.ComponentBaseDefinition {
	return zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"lighthouseAthena": GetConsensusClientSkeletonBase(network),
		},
	}
}

func GetConsensusClientSkeletonBase(network string) zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition {
	args := []string{"-c", "/scripts/lighthouse.sh"}
	switch network {
	case hestia_req_types.Goerli:
		args = []string{"-c", "/scripts/lighthouseGoerliBeacon.sh"}
	}
	sbCfg := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ConsensusClientChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			StatefulSetDriver: &zeus_topology_config_drivers.StatefulSetDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					athena: {
						Container: v1Core.Container{
							Env: ClusterConfigEnvVars(nil, network),
						},
					},
					consensusClient: {
						Container: v1Core.Container{
							Name:  consensusClient,
							Image: lighthouseDockerImage,
							Args:  args,
							Env:   ClusterConfigEnvVars(nil, network),
						},
					},
					initSnapshots: {Container: v1Core.Container{
						Name: initSnapshots,
						Env:  ClusterConfigEnvVars(nil, network),
					}},
				},
			},
		},
	}
	return sbCfg
}

func GetExecClientComponentBaseConfig(network string) zeus_cluster_config_drivers.ComponentBaseDefinition {
	return zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"gethAthena": GetExecClientSkeletonBase(network),
		},
	}
}

func GetExecClientSkeletonBase(network string) zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition {
	args := []string{"-c", "/scripts/geth.sh"}
	switch network {
	case hestia_req_types.Goerli:
		args = []string{"-c", "/scripts/gethGoerli.sh"}
	}
	sbCfg := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ExecClientChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			StatefulSetDriver: &zeus_topology_config_drivers.StatefulSetDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					athena: {
						Container: v1Core.Container{
							Env: ClusterConfigEnvVars(nil, network),
						},
					},
					execClient: {
						Container: v1Core.Container{
							Name:  execClient,
							Image: gethDockerImage,
							Args:  args,
							Env:   ClusterConfigEnvVars(nil, network),
						},
					},
					initSnapshots: {Container: v1Core.Container{
						Name: initSnapshots,
						Env:  ClusterConfigEnvVars(nil, network),
					}},
				},
			},
		},
	}
	return sbCfg
}
