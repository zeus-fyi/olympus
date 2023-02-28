package olympus_beacon_cookbooks

import (
	olympus_hydra_choreography_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra/choreography"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1Core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	consensusClient       = "zeus-consensus-client"
	execClient            = "zeus-exec-client"
	lighthouseDockerImage = "sigp/lighthouse:v3.5.0"
	gethDockerImage       = "ethereum/client-go:v1.11.2"
	initSnapshots         = "init-snapshots"
)

var (
	EphemeralBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumEphemeralAthenaBeacon",
		CloudCtxNs:       EphemeralAthenaBeaconCloudCtxNs,
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  ConsensusClientComponentBase,
			"execClients":       ExecClientComponentBase,
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	MainnetBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumMainnetAthenaBeacon",
		CloudCtxNs:       MainnetAthenaBeaconCloudCtxNs,
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  ConsensusClientComponentBaseMainnet,
			"execClients":       ExecClientComponentBaseMainnet,
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	ConsensusClientComponentBaseMainnet = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"lighthouseAthena": ConsensusClientSkeletonBaseConfigMainnet,
		},
	}
	ConsensusClientSkeletonBaseConfigMainnet = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ConsensusClientChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			ConfigMapDriver: &zeus_topology_config_drivers.ConfigMapDriver{
				ConfigMap: v1Core.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "cm-consensus-client"},
				},
			},
			StatefulSetDriver: &zeus_topology_config_drivers.StatefulSetDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					consensusClient: {Container: v1Core.Container{
						Name:  consensusClient,
						Image: lighthouseDockerImage,
						Args:  []string{"-c", "/scripts/lighthouse.sh"},
						Env:   ClusterConfigEnvVars(nil, "mainnet"),
					}},
					initSnapshots: {Container: v1Core.Container{
						Name: initSnapshots,
						Env:  ClusterConfigEnvVars(nil, "mainnet"),
					}},
				},
			},
		}}

	MainnetAthenaBeaconCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "athena-beacon-mainnet", // set with your own namespace
		Env:           "production",
	}
	EphemeralAthenaBeaconCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "athena-beacon-ephemeral", // set with your own namespace
		Env:           "production",
	}
	ConsensusClientComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"lighthouseAthena": ConsensusClientSkeletonBaseConfig,
		},
	}
	ExecClientComponentBaseMainnet = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"gethAthena": ExecClientSkeletonBaseConfigMainnet,
		},
	}
	ExecClientSkeletonBaseConfigMainnet = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ExecClientChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			ConfigMapDriver: &zeus_topology_config_drivers.ConfigMapDriver{
				ConfigMap: v1Core.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "cm-exec-client"},
				},
			},
			StatefulSetDriver: &zeus_topology_config_drivers.StatefulSetDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					execClient: {
						Container: v1Core.Container{
							Name:  execClient,
							Image: gethDockerImage,
							Args:  []string{"-c", "/scripts/geth.sh"},
							Env:   ClusterConfigEnvVars(nil, "mainnet"),
						},
					},
					initSnapshots: {Container: v1Core.Container{
						Name: initSnapshots,
						Env:  ClusterConfigEnvVars(nil, "mainnet"),
					}},
				},
			},
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
