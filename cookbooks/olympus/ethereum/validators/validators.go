package olympus_hydra_validators_cookbooks

import (
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
)

const HydraValidatorsClientName = "zeus-hydra-validators"

var (
	ValidatorSkeletonBaseName          = "lighthouseAthenaValidatorClient"
	ValidatorSecondarySkeletonBaseName = "lighthouseAthenaValidatorClientSecondary"

	ValidatorClientComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			ValidatorSkeletonBaseName: ValidatorClientClientSkeletonBaseConfig,
		},
	}
	ValidatorClientSecondaryComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			ValidatorSecondarySkeletonBaseName: ValidatorClientClientSecondarySkeletonBaseConfig,
		},
	}
	ValidatorClientClientSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: ValidatorsChartPath,
	}
	ValidatorsChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/validators/infra",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        ValidatorSkeletonBaseName, // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
	ValidatorClientClientSecondarySkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: ValidatorsSecondaryChartPath,
	}
	ValidatorsSecondaryChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/validators/secondary_infra",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        ValidatorSkeletonBaseName + "secondary", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)
