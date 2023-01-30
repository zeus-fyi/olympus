package olympus_hydra_cookbooks

import (
	"fmt"

	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
)

const (
	protocolNetworkKeyEnv = "PROTOCOL_NETWORK_ID"
	ephemeryNamespace     = "ephemeral-staking"
	mainnetNamespace      = "mainnet-staking"
)

var (
	HydraPort = v1.ContainerPort{
		Name:          "hydra",
		ContainerPort: 9000,
		Protocol:      v1.Protocol("TCP"),
	}
	HydraContainer = zeus_topology_config_drivers.ContainerDriver{
		Container: v1.Container{
			Name:            "hydra",
			Image:           "registry.digitalocean.com/zeus-fyi/hydra:latest",
			Ports:           []v1.ContainerPort{HydraPort},
			ImagePullPolicy: "Always",
		}}
)

func HydraClusterConfig(network string) zeus_cluster_config_drivers.ClusterDefinition {
	cd := HydraClusterDefinition
	switch network {
	case "mainnet":
		cd.CloudCtxNs.Namespace = mainnetNamespace
		cd.ClusterClassName = "hydraMainnet"

		depCfgOverride := zeus_topology_config_drivers.DeploymentDriver{}
		depCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
		depCfgOverride.ContainerDrivers["hydra"] = HydraContainerDriver(network)
		cfgOverride := zeus_topology_config_drivers.TopologyConfigDriver{
			IngressDriver:     nil,
			StatefulSetDriver: nil,
			ServiceDriver:     nil,
			DeploymentDriver:  &depCfgOverride,
		}
		tmp := HydraComponentBase
		sb := tmp.SkeletonBases["hydra"]
		sb.TopologyConfigDriver = &cfgOverride
		tmp.SkeletonBases["hydra"] = sb
		cd.ComponentBases["hydra"] = tmp
	case "ephemery":
		cd.CloudCtxNs.Namespace = ephemeryNamespace
		cd.ClusterClassName = "hydraEphemery"

		depCfgOverride := zeus_topology_config_drivers.DeploymentDriver{}
		depCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
		depCfgOverride.ContainerDrivers["hydra"] = HydraContainerDriver(network)
		cfgOverride := zeus_topology_config_drivers.TopologyConfigDriver{
			IngressDriver:     nil,
			StatefulSetDriver: nil,
			ServiceDriver:     nil,
			DeploymentDriver:  &depCfgOverride,
		}
		tmp := HydraComponentBase
		sb := tmp.SkeletonBases["hydra"]
		sb.TopologyConfigDriver = &cfgOverride
		tmp.SkeletonBases["hydra"] = sb
		cd.ComponentBases["hydra"] = tmp
	}
	return cd
}

func HydraContainerDriver(network string) zeus_topology_config_drivers.ContainerDriver {
	switch network {
	case "mainnet":
		envVar := HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
		HydraContainer.AppendEnvVars = []v1.EnvVar{envVar}
	case "ephemery":
		envVar := HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))
		HydraContainer.AppendEnvVars = []v1.EnvVar{envVar}
	}
	return HydraContainer
}
