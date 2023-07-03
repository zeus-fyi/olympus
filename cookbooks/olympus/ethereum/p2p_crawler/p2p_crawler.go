package p2p_crawler

import (
	"fmt"
	"time"

	olympus_hydra_choreography_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra/choreography"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	P2PCrawlerClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "p2pCrawler",
		CloudCtxNs:       P2PCloudCtxNs,
		ComponentBases:   BeaconComponentBases,
	}
	BeaconComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"p2pCrawler":        P2PCrawlerComponentBase,
		"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
	}
	P2PCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "p2p-crawler", // set with your own namespace
		Env:           "production",
	}
	P2PCrawlerComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"p2pCrawler": P2PCrawlerSkeletonBaseConfig,
		},
	}
	P2PCrawlerSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: P2PCrawlerChartPath,
	}
	P2PCrawlerUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "p2pCrawler",
		ChartName:         "p2pCrawler",
		ChartDescription:  "p2pCrawler",
		SkeletonBaseName:  "p2pCrawler",
		ComponentBaseName: "p2pCrawler",
		ClusterClassName:  "p2pCrawler",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	P2PCrawlerChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/p2p_crawler/infra",
		DirOut:      "./olympus/outputs",
		FnIn:        "p2pCrawler", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
	}
)
