package olympus_hydra_choreography_cookbooks

import (
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
)

var (
	HydraChoreographyComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"hydraChoreography": HydraChoreographySkeletonBaseConfig,
		},
	}
	HydraChoreographySkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: HydraChoreographyChartPath,
	}
	HydraChoreographyChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/hydra/choreography/infra",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "hydraChoreography", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)
