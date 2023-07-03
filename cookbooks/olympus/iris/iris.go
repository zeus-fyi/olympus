package iris_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	irisClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "iris",
		CloudCtxNs:       irisCloudCtxNs,
		ComponentBases:   irisComponentBases,
	}
	irisComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"iris": irisComponentBase,
	}
	irisComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"iris": irisSkeletonBaseConfig,
		},
	}
	irisSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: irisChartPath,
	}
)

var irisUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:      "iris",
	ChartName:         "iris",
	ChartDescription:  "iris",
	SkeletonBaseName:  "iris",
	ComponentBaseName: "iris",
	ClusterClassName:  "iris",
	Tag:               "latest",
	Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var irisCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "iris", // set with your own namespace
	Env:           "production",
}

var irisDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: irisCloudCtxNs,
}

var irisChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/iris/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "iris", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
}
