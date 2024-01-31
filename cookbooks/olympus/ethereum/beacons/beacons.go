package olympus_beacon_cookbooks

import (
	"fmt"

	olympus_hydra_choreography_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra/choreography"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/zeus/workload_config_drivers/config_overrides"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
	v1Core "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	athena                = "athena"
	consensusClient       = "zeus-consensus-client"
	execClient            = "zeus-exec-client"
	lighthouseDockerImage = "sigp/lighthouse:v3.5.1"
	gethDockerImage       = "ethereum/client-go:v1.11.5"
	initSnapshots         = "init-snapshots"
)

var (
	EphemeralBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumEphemeralAthenaBeacon",
		CloudCtxNs:       GetBeaconCloudCtxNs(hestia_req_types.Ephemery),
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  GetConsensusClientComponentBaseConfig(hestia_req_types.Ephemery),
			"execClients":       GetExecClientComponentBaseConfig(hestia_req_types.Ephemery),
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	GoerliBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumGoerliAthenaBeacon",
		CloudCtxNs:       GetBeaconCloudCtxNs(hestia_req_types.Goerli),
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  GetConsensusClientComponentBaseConfig(hestia_req_types.Goerli),
			"execClients":       GetExecClientComponentBaseConfig(hestia_req_types.Goerli),
			"beaconIngress":     GetIngressComponentBaseConfig(hestia_req_types.Goerli),
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	MainnetBeaconBaseClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "ethereumMainnetAthenaBeacon",
		CloudCtxNs:       GetBeaconCloudCtxNs(hestia_req_types.Mainnet),
		ComponentBases: map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
			"consensusClients":  GetConsensusClientComponentBaseConfig(hestia_req_types.Mainnet),
			"execClients":       GetExecClientComponentBaseConfig(hestia_req_types.Mainnet),
			"hydraChoreography": olympus_hydra_choreography_cookbooks.HydraChoreographyComponentBase,
		},
	}
	ConsensusClientComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"lighthouseAthena": ConsensusClientSkeletonBaseConfig,
		},
	}
	ConsensusClientSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ConsensusClientChartPath,
	}
	ConsensusClientChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/beacons/infra/consensus_client",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "lighthouseAthena", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
	ExecClientComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"gethAthena": ExecClientSkeletonBaseConfig,
		},
	}
	ExecClientSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ExecClientChartPath,
	}
	ExecClientChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/beacons/infra/exec_client",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "gethAthena", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
	IngressChartPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./olympus/ethereum/beacons/infra/ingress",
		DirOut:      "./olympus/ethereum/outputs",
		FnIn:        "ingress", // filename for your gzip workload
		FnOut:       "",
		Env:         "",
	}
)

func GetBeaconCloudCtxNs(network string) zeus_common_types.CloudCtxNs {
	beaconCloudCtxNs := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     fmt.Sprintf("athena-beacon-%s", network), // set with your own namespace
		Env:           "production",
	}
	return beaconCloudCtxNs
}

func GetConsensusClientComponentBaseConfig(network string) zeus_cluster_config_drivers.ComponentBaseDefinition {
	return zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"lighthouseAthena": GetConsensusClientSkeletonBase(network),
		},
	}
}

func GetConsensusClientSkeletonBase(network string) zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition {
	args := []string{"-c", "/scripts/lighthouse.sh"}
	switch network {
	case hestia_req_types.Goerli:
		args = []string{"-c", "/scripts/lighthouseGoerliBeacon.sh"}
	}
	sd := &zeus_topology_config_drivers.ServiceDriver{}
	rrCC := v1Core.ResourceRequirements{}

	switch network {
	case hestia_req_types.Goerli:
		rrCC = v1Core.ResourceRequirements{
			Limits: v1Core.ResourceList{
				"cpu":    resource.MustParse(consensusClientGoerliRequestLimitCPU),
				"memory": resource.MustParse(consensusClientGoerliRequestLimitRAM),
			},
			Requests: v1Core.ResourceList{
				"cpu":    resource.MustParse(consensusClientGoerliRequestCPU),
				"memory": resource.MustParse(consensusClientGoerliRequestRAM),
			},
		}
	}
	sd.AddNginxTargetPort("http", "http-api")
	sbCfg := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ConsensusClientChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			ServiceDriver: sd,
			StatefulSetDriver: &zeus_topology_config_drivers.StatefulSetDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					athena: {
						Container: v1Core.Container{
							Env: ClusterConfigEnvVars(nil, network),
						},
					},
					consensusClient: {
						Container: v1Core.Container{
							Name:      consensusClient,
							Image:     lighthouseDockerImage,
							Args:      args,
							Env:       ClusterConfigEnvVars(nil, network),
							Resources: rrCC,
						},
					},
					initSnapshots: {Container: v1Core.Container{
						Name: initSnapshots,
						Env:  ClusterConfigEnvVars(nil, network),
					}},
				},
			},
		},
	}
	return sbCfg
}

func GetExecClientComponentBaseConfig(network string) zeus_cluster_config_drivers.ComponentBaseDefinition {
	return zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"gethAthena": GetExecClientSkeletonBase(network),
		},
	}
}

func GetExecClientSkeletonBase(network string) zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition {
	args := []string{"-c", "/scripts/geth.sh"}
	switch network {
	case hestia_req_types.Goerli:
		args = []string{"-c", "/scripts/gethGoerli.sh"}
	}
	rrEC := v1Core.ResourceRequirements{}
	switch network {
	case hestia_req_types.Goerli:
		rrEC = v1Core.ResourceRequirements{
			Limits: v1Core.ResourceList{
				"cpu":    resource.MustParse(execClientGoerliRequestLimitCPU),
				"memory": resource.MustParse(execClientGoerliRequestLimitRAM),
			},
			Requests: v1Core.ResourceList{
				"cpu":    resource.MustParse(execClientGoerliRequestCPU),
				"memory": resource.MustParse(execClientGoerliRequestRAM),
			},
		}
	}

	sd := &zeus_topology_config_drivers.ServiceDriver{}
	sd.AddNginxTargetPort("http", "http-rpc")
	sbCfg := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: ExecClientChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			ServiceDriver: sd,
			StatefulSetDriver: &zeus_topology_config_drivers.StatefulSetDriver{
				ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
					athena: {
						Container: v1Core.Container{
							Env: ClusterConfigEnvVars(nil, network),
						},
					},
					execClient: {
						Container: v1Core.Container{
							Name:      execClient,
							Image:     gethDockerImage,
							Args:      args,
							Env:       ClusterConfigEnvVars(nil, network),
							Resources: rrEC,
						},
					},
					initSnapshots: {Container: v1Core.Container{
						Name: initSnapshots,
						Env:  ClusterConfigEnvVars(nil, network),
					}},
				},
			},
		},
	}
	return sbCfg
}

func GetIngressComponentBaseConfig(network string) zeus_cluster_config_drivers.ComponentBaseDefinition {
	return zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"ingress": GetIngressSkeletonBase(network),
		},
	}
}
func GetIngressSkeletonBase(network string) zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition {
	sbCfg := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseChart:         zeus_req_types.TopologyCreateRequest{},
		SkeletonBaseNameChartPath: IngressChartPath,
		TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{
			IngressDriver: &zeus_topology_config_drivers.IngressDriver{
				Ingress: v1.Ingress{
					Spec: v1.IngressSpec{
						TLS: []v1.IngressTLS{{[]string{fmt.Sprintf("eth.%s.zeus.fyi", network)}, fmt.Sprintf("beacon-%s-tls", network)}},
					},
				},
				Host:         fmt.Sprintf("eth.%s.zeus.fyi", network),
				NginxAuthURL: "https://aegis.zeus.fyi/v1beta/ethereum/beacon",
			},
		},
	}
	return sbCfg
}
