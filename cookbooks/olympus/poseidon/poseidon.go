package poseidon_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	PoseidonClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "poseidon",
		CloudCtxNs:       PoseidonCloudCtxNs,
		ComponentBases:   PoseidonComponentBases,
	}
	PoseidonComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"poseidon": PoseidonComponentBase,
	}
	PoseidonComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"poseidon": PoseidonSkeletonBaseConfig,
		},
	}
	PoseidonSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: PoseidonChartPath,
	}

	PoseidonCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "ovh",
		Region:        "us-west-or-1",
		Context:       "kubernetes-admin@zeusfyi",
		Namespace:     "poseidon", // set with your own namespace
		Env:           "production",
	}
	PoseidonDeployKnsReq = zeus_req_types.TopologyDeployRequest{
		TopologyID: 0,
		CloudCtxNs: PoseidonCloudCtxNs,
	}
	PoseidonChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/poseidon/infra",
		DirOut:      "./olympus/outputs",
		FnIn:        "poseidon", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)

var PoseidonUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "poseidon",
	ChartName:        "poseidon",
	ChartDescription: "poseidon",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}
