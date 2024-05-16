package zeus_cookbook

import (
	"context"
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var ctx = context.Background()

var (
	flowsClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "flows",
		ComponentBases:   flowsComponentBases,
	}
	flowsComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"flows": flowsComponentBase,
	}
	flowsComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"flows": flowsSkeletonBaseConfig,
		},
	}
	flowsSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: flowsChartPath,
	}
)

// chart workload metadata
var uploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "flows",
	ChartName:        "flows",
	ChartDescription: "flows",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var FlowsCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "flows", // set with your own namespace
	Env:           "production",
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var flowsChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/flows/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "flows", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
}
