package hestia_olympus_cookbook

import (
	"fmt"
	"time"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	clusterClassName    = "olympus"
	ZeusCloudClusterDef = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: clusterClassName,
		ComponentBases:   ZeusCloudComponentBases,
		CloudCtxNs:       ZeusCloudCloudCtxNs,
	}
	ZeusCloudComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"zeusCloud": ZeusCloudComponentBase,
	}
	ZeusCloudComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"zeusCloud": ZeusCloudSkeletonBaseConfig,
		},
	}
	ZeusCloudCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "hestia", // set with your own namespace
		Env:           "production",
	}
	ZeusCloudSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ZeusCloudChartPath,
	}
	ZeusCloudUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "zeusCloud",
		ChartName:         "zeusCloud",
		ChartDescription:  "zeusCloud",
		SkeletonBaseName:  "zeusCloud",
		ComponentBaseName: "zeusCloud",
		ClusterClassName:  "olympus",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	ZeusCloudChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/hestia/frontend_infra",
		DirOut:      "./olympus/outputs",
		FnIn:        "zeusCloud", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
	}
)

var ZeusCloudDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: ZeusCloudCloudCtxNs,
}
