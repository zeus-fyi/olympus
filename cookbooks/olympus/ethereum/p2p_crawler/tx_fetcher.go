package p2p_crawler

import (
	"fmt"
	"time"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
)

var (
	TxFetcherClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "txFetcher",
		CloudCtxNs:       TxFetcherCloudCtxNs,
		ComponentBases:   TxFetcherComponentBases,
	}
	TxFetcherComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"txFetcher": TxFetcherComponentBase,
	}
	TxFetcherCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "tx-fetcher",
		Env:           "production",
	}
	TxFetcherComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"txFetcher": TxFetcherSkeletonBaseConfig,
		},
	}
	TxFetcherSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: TxFetcherChartPath,
	}
	TxFetcherUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "txFetcher",
		ChartName:         "txFetcher",
		ChartDescription:  "txFetcher",
		SkeletonBaseName:  "txFetcher",
		ComponentBaseName: "txFetcher",
		ClusterClassName:  "txFetcher",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	TxFetcherChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/p2p_crawler/txfetcher",
		DirOut:      "./olympus/outputs",
		FnIn:        "txFetcher", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
	}
)
