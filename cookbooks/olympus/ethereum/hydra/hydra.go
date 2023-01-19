package olympus_hydra_cookbooks

import (
	olympus_beacon_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/beacons"
	olympus_hydra_choreography_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra/choreography"
	olympus_hydra_validators_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/validators"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
)

var (
	EphemeralHydraClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "hydraEphemeral",
		CloudCtxNs:       ValidatorCloudCtxNs,
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients": olympus_beacon_cookbooks.ConsensusClientComponentBase,
			"execClients":      olympus_beacon_cookbooks.ExecClientComponentBase,
			"validatorClients": olympus_hydra_validators_cookbooks.ValidatorClientComponentBase,
			"choreography":     olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
			"hydra":            HydraComponentBase,
		},
	}
	ValidatorCloudCtxNs = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking", // set with your own namespace
		Env:           "production",
	}
	HydraComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"hydra": HydraSkeletonBaseConfig,
		},
	}
	HydraSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: HydraChartPath,
	}
	HydraChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/hydra/remote_signer/infra",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "hydra", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)

// TODO config workload types and env vars, and for cli inputs to athena, and snapshotters
/*
	Cmd.Flags().StringVar(&bearer, "bearer", "", "bearer for cluster actions")
	Cmd.Flags().IntVar(&AthenaProtocolCloudCtxNs.ProtocolNetworkID, "protocol-network-id", 0, "identifier for protocol and network")
	Cmd.Flags().StringVar(&AthenaProtocolCloudCtxNs.CloudCtxNs.CloudProvider, "cloud-provider", "", "cloud-provider")
	Cmd.Flags().StringVar(&AthenaProtocolCloudCtxNs.CloudCtxNs.Context, "ctx", "", "context")
	Cmd.Flags().StringVar(&AthenaProtocolCloudCtxNs.CloudCtxNs.Namespace, "ns", "", "namespace")
	Cmd.Flags().StringVar(&AthenaProtocolCloudCtxNs.CloudCtxNs.Region, "region", "", "region")
*/
