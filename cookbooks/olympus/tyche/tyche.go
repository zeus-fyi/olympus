package tyche_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	TycheClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "tyche",
		CloudCtxNs:       tycheCloudCtxNs,
		ComponentBases:   tycheComponentBases,
	}
	tycheComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"tyche": tycheComponentBase,
	}
	tycheComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"tyche": tycheSkeletonBaseConfig,
		},
	}
	tycheSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: tycheChartPath,
	}

	tycheCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "ovh",
		Region:        "us-west-or-1",
		Context:       "kubernetes-admin@zeusfyi",
		Namespace:     "tyche", // set with your own namespace
		Env:           "production",
	}
	tycheDeployKnsReq = zeus_req_types.TopologyDeployRequest{
		TopologyID: 0,
		CloudCtxNs: tycheCloudCtxNs,
	}
	tycheChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/tyche/infra",
		DirOut:      "./olympus/outputs",
		FnIn:        "tyche", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)

var tycheUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "tyche",
	ChartName:        "tyche",
	ChartDescription: "tyche",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}
