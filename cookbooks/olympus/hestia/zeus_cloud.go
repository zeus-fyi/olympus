package hestia_olympus_cookbook

import (
	"fmt"
	"time"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

var (
	clusterClassName    = "info-flows"
	ZeusCloudClusterDef = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: clusterClassName,
		ComponentBases:   ZeusCloudComponentBases,
		CloudCtxNs:       ZeusCloudCloudCtxNs,
	}
	ZeusCloudComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"info-flows": ZeusCloudComponentBase,
	}
	ZeusCloudComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"info-flows": ZeusCloudSkeletonBaseConfig,
		},
	}
	clusterClassNameStaging     = "info-flows-staging"
	FlowsStagingCloudClusterDef = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: clusterClassNameStaging,
		ComponentBases:   FlowsStagingComponentBases,
	}
	FlowsStagingComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"info-flows-staging": FlowsStagingComponentBase,
	}
	FlowsStagingComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"info-flows-staging": FlowsStagingSkeletonBaseConfig,
		},
	}
	FlowsStagingSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: FlowsStagingChartPath,
	}
	FlowsStagingChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/hestia/info_works_staging",
		DirOut:      "./olympus/outputs",
		FnIn:        "info-flows-staging", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
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
		TopologyName:      "info-flows",
		ChartName:         "info-flows",
		ChartDescription:  "info-flows",
		SkeletonBaseName:  "info-flows",
		ComponentBaseName: "info-flows",
		ClusterClassName:  "info-flows",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	ZeusCloudChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/hestia/info_works",
		DirOut:      "./olympus/outputs",
		FnIn:        "info-flows", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
	}
)

var ZeusCloudDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: ZeusCloudCloudCtxNs,
}
