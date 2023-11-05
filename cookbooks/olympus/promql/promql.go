package olympus_cookbooks_promql

import (
	"fmt"
	"time"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

var (
	promqlClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "promql",
		ComponentBases:   promqlComponentBases,
	}
	promqlComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"promql": promqlComponentBase,
	}
	promqlComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"promql": promqlSkeletonBaseConfig,
		},
	}
	promqlSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: promqlChartPath,
	}
)

var (
	promqlUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "promql",
		ChartName:         "promql",
		ChartDescription:  "promql",
		SkeletonBaseName:  "promql",
		ComponentBaseName: "promql",
		ClusterClassName:  "promql",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	promqlChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/promql/infra",
		DirOut:      "./olympus/outputs",
		FnIn:        "promql", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)
