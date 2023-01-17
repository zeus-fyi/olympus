package olympus_beacon_cookbooks

import (
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
)

var (
	BeaconComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"consensusClients": ConsensusClientComponentBase,
		"execClients":      ExecClientComponentBase,
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
		DirIn:       "./olympus/ethereum/beacons/infra/consensus_client",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "gethAthena", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)
