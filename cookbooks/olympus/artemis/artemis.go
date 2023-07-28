package artemis_cookbook

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
	ArtemisClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "artemis",
		CloudCtxNs:       ArtemisCloudCtxNs,
		ComponentBases:   BeaconComponentBases,
	}
	BeaconComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"artemis": ArtemisComponentBase,
	}
	ArtemisCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "ovh",
		Region:        "us-west-or-1",
		Context:       "kubernetes-admin@zeusfyi",
		Namespace:     "artemis", // set with your own namespace
		Env:           "production",
	}
	ArtemisComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"artemis": ArtemisSkeletonBaseConfig,
		},
	}
	ArtemisSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: ArtemisUploadChartPath,
	}
	ArtemisUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "artemis",
		ChartName:         "artemis",
		ChartDescription:  "artemis",
		SkeletonBaseName:  "artemis",
		ComponentBaseName: "artemis",
		ClusterClassName:  "artemis",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	ArtemisUploadChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/artemis/infra",
		DirOut:      "./olympus/outputs",
		FnIn:        "artemis", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)
