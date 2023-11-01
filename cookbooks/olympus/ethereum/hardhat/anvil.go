package olympus_hardhat

import (
	"fmt"
	"time"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

/*
https://book.getfoundry.sh/reference/anvil/
https://book.getfoundry.sh/tutorials/forking-mainnet-with-cast-anvil
anvil --fork-url https://mainnet.infura.io/v3/$INFURA_KEY

anvil_autoImpersonateAccount
Accepts true to enable auto impersonation of accounts, and false to disable it. When enabled, any transaction's sender will be automatically impersonated. Same as anvil_impersonateAccount.

anvil_reset
Reset the fork to a fresh forked state, and optionally update the fork config


docker tag ghcr.io/foundry-rs/foundry:latest foundry:latest
*/

var (
	anvilClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "anvil",
		CloudCtxNs:       anvilCtxNs,
		ComponentBases:   anvilComponentBases,
	}
	anvilCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "anvil", // set with your own namespace
		Env:           "production",
	}
	anvilComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"anvil": anvilComponentBase,
	}
	anvilComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"anvil": anvilSkeletonBaseConfig,
		},
	}
	anvilSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: anvilChartPath,
	}
	anvilUploadChart = zeus_req_types.TopologyCreateRequest{
		TopologyName:      "anvil",
		ChartName:         "anvil",
		ChartDescription:  "anvil",
		SkeletonBaseName:  "anvil",
		ComponentBaseName: "anvil",
		ClusterClassName:  "anvil",
		Tag:               "latest",
		Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
	anvilChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/hardhat/anvil",
		DirOut:      "./olympus/outputs",
		FnIn:        "anvil", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
	}

	serverlessAnvilClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "anvil-serverless",
		CloudCtxNs:       anvilCtxNs,
		ComponentBases:   serverlessAnvilComponentBases,
	}
	serverlessAnvilComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"anvil-serverless": serverlessAnvilComponentBase,
	}
	serverlessAnvilComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"anvil-serverless": serverlessAnvilSkeletonBaseConfig,
		},
	}
	serverlessAnvilSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: anvilServerlessChartPath,
	}
	anvilServerlessChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/hardhat/serverless_anvil",
		DirOut:      "./olympus/outputs",
		FnIn:        "anvil-serverless", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
		FilterFiles: &strings_filter.FilterOpts{},
	}
)
