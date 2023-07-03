package olympus_hardhat

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
	hardhatClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "hardhat",
		CloudCtxNs:       hardhatCtxNs,
		ComponentBases:   HardhatComponentBases,
	}
	hardhatCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "hardhat", // set with your own namespace
		Env:           "production",
	}
	HardhatComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"hardhat": hardhatComponentBase,
	}
	hardhatComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"hardhat": hardhatSkeletonBaseConfig,
		},
	}
	hardhatSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: hardhatChartPath,
	}
	hardhatUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "hardhat",
		ChartName:         "hardhat",
		ChartDescription:  "hardhat",
		SkeletonBaseName:  "hardhat",
		ComponentBaseName: "hardhat",
		ClusterClassName:  "hardhat",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	hardhatChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/hardhat/hardhat",
		DirOut:      "./olympus/outputs",
		FnIn:        "hardhat", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
	}
)
