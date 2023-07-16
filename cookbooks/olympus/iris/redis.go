package iris_olympus_cookbook

import (
	"fmt"
	"time"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

var (
	redisClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "redis",
		CloudCtxNs:       redisCloudCtxNs,
		ComponentBases:   redisComponentBases,
	}
	// todo: add replicas
	redisComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"master":   redisComponentBase,
		"replicas": redisComponentBase,
	}
	redisComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"redis": redisSeletonBaseConfig,
		},
	}
	redisSeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: redisChartPath,
	}
)

var redisUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:      "redis",
	ChartName:         "redis",
	ChartDescription:  "redis",
	SkeletonBaseName:  "redis",
	ComponentBaseName: "redis",
	ClusterClassName:  "redis",
	Tag:               "latest",
	Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var redisCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "redis", // set with your own namespace
	Env:           "production",
}

var redisDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: redisCloudCtxNs,
}

var redisChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/iris/redis/master",
	DirOut:      "./olympus/outputs",
	FnIn:        "redis", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
}
