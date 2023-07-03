package aegis_olympus_cookbook

import (
	"fmt"
	"time"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	AegisClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "aegis",
		CloudCtxNs:       AegisCloudCtxNs,
		ComponentBases:   AegisComponentBases,
	}
	AegisComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"aegis": AegisComponentBase,
	}
	AegisComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"aegis": AegisSkeletonBaseConfig,
		},
	}
	AegisSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: AegisChartPath,
	}
	AegisUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "aegis",
		ChartName:         "aegis",
		ChartDescription:  "aegis",
		SkeletonBaseName:  "aegis",
		ComponentBaseName: "aegis",
		ClusterClassName:  "aegis",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	AegisCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "ovh",
		Region:        "us-west-or-1",
		Context:       "kubernetes-admin@zeusfyi",
		Namespace:     "aegis", // set with your own namespace
		Env:           "production",
	}
	AegisDeployKnsReq = zeus_req_types.TopologyDeployRequest{
		TopologyID: 1669423376281749000,
		CloudCtxNs: AegisCloudCtxNs,
	}
	AegisChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/aegis/infra",
		DirOut:      "./olympus/outputs",
		FnIn:        "aegis", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)
