package olympus_ethereum_mev_cookbooks

import (
	"context"

	ethereum_mev_cookbooks "github.com/zeus-fyi/zeus/cookbooks/ethereum/mev"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/zeus/workload_config_drivers/config_overrides"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	v1Core "k8s.io/api/core/v1"
)

const (
	mevContainerReference = "zeus-mev"
	flashbotsDockerImage  = "flashbots/mev-boost:1.5.0"
)

var (
	ctx                   = context.Background()
	MevSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: MevChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			DeploymentDriver: &zeus_topology_config_drivers.DeploymentDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					mevContainerReference: {Container: v1Core.Container{
						Name:  mevContainerReference,
						Image: flashbotsDockerImage,
						Args:  ethereum_mev_cookbooks.GetMevBoostArgs(ctx, hestia_req_types.Goerli, MevRelays),
					}},
				},
			}},
	}
	MevChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/mev/infra",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "mev", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
	MevRelays = ethereum_mev_cookbooks.RelaysEnabled{
		Flashbots:   true,
		Blocknative: true,
		EdenNetwork: true,
	}
)

func MevCbCfgGoerli() zeus_cluster_config_drivers.ComponentBaseDefinition {
	MevSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: MevChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			DeploymentDriver: &zeus_topology_config_drivers.DeploymentDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					mevContainerReference: {Container: v1Core.Container{
						Name:  mevContainerReference,
						Image: flashbotsDockerImage,
						Args:  ethereum_mev_cookbooks.GetMevBoostArgs(ctx, hestia_req_types.Goerli, MevRelays),
					}},
				},
			}},
	}
	MevComponentBase := zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"mevBoost": MevSkeletonBaseConfig,
		},
	}
	return MevComponentBase
}

func MevCbCfgMainnet() zeus_cluster_config_drivers.ComponentBaseDefinition {
	MevSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: MevChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			DeploymentDriver: &zeus_topology_config_drivers.DeploymentDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					mevContainerReference: {Container: v1Core.Container{
						Name:  mevContainerReference,
						Image: flashbotsDockerImage,
						Args:  ethereum_mev_cookbooks.GetMevBoostArgs(ctx, hestia_req_types.Mainnet, MevRelays),
					}},
				},
			}},
	}
	MevComponentBase := zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"mevBoost": MevSkeletonBaseConfig,
		},
	}
	return MevComponentBase
}
