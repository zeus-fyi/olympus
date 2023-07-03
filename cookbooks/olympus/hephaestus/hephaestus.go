package hephaestus_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	HephaestusClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "hephaestus",
		CloudCtxNs:       HephaestusCloudCtxNs,
		ComponentBases:   HephaestusComponentBases,
	}
	HephaestusComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"hephaestus": HephaestusComponentBase,
	}
	HephaestusComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"hephaestus": HephaestusSkeletonBaseConfig,
		},
	}
	HephaestusSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: HephaestusChartPath,
	}
)

var HephaestusUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:      "hephaestus",
	ChartName:         "hephaestus",
	ChartDescription:  "hephaestus",
	SkeletonBaseName:  "hephaestus",
	ComponentBaseName: "hephaestus",
	ClusterClassName:  "hephaestus",
	Tag:               "latest",
	Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var HephaestusCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "hephaestus", // set with your own namespace
	Env:           "production",
}

var HephaestusDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: HephaestusCloudCtxNs,
}

var HephaestusChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/hephaestus/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "hephaestus", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
}
