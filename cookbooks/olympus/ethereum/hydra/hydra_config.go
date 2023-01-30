package olympus_hydra_cookbooks

import (
	"fmt"

	olympus_common_vals_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/common"
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

func HydraClusterConfig(cd *zeus_cluster_config_drivers.ClusterDefinition, network string) *zeus_cluster_config_drivers.ClusterDefinition {

	var envVar v1.EnvVar
	switch network {
	case "mainnet":
		cd.CloudCtxNs.Namespace = mainnetNamespace
		cd.ClusterClassName = "hydraMainnet"

		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
	case "ephemery":
		cd.CloudCtxNs.Namespace = ephemeryNamespace
		cd.ClusterClassName = "hydraEphemery"
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))
	}

	envVarsChoreography := olympus_common_vals_cookbooks.GetChoreographyEnvVars()
	internalAuthEnvVars := olympus_common_vals_cookbooks.GetCommonInternalAuthEnvVars()
	combinedEnvVars := append(envVarsChoreography, internalAuthEnvVars...)

	containCfg := zeus_topology_config_drivers.ContainerDriver{}
	containCfg.Env = envVarsChoreography
	containCfg.AppendEnvVars = []v1.EnvVar{envVar}

	containCfgBoth := zeus_topology_config_drivers.ContainerDriver{}
	containCfgBoth.Env = combinedEnvVars
	containCfgBoth.AppendEnvVars = []v1.EnvVar{envVar}

	depCfgOverride := zeus_topology_config_drivers.DeploymentDriver{}
	depCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)

	depCfgOverride.ContainerDrivers["hydra"] = containCfgBoth
	depCfgOverride.ContainerDrivers["zeus-hydra-choreography"] = containCfg
	depCfgOverride.ContainerDrivers["athena"] = containCfgBoth

	stsCfgOverride := zeus_topology_config_drivers.StatefulSetDriver{}
	stsCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	stsCfgOverride.ContainerDrivers["athena"] = containCfgBoth
	stsCfgOverride.ContainerDrivers["zeus-consensus-client"] = containCfgBoth
	stsCfgOverride.ContainerDrivers["zeus-exec-client"] = containCfgBoth
	stsCfgOverride.ContainerDrivers["init-validator"] = containCfgBoth
	stsCfgOverride.ContainerDrivers["init-snapshot"] = containCfgBoth

	for k, v := range cd.ComponentBases {
		if k == "hydra" || k == "hydraChoreography" {
			cfgOverride := zeus_topology_config_drivers.TopologyConfigDriver{
				IngressDriver:     nil,
				StatefulSetDriver: nil,
				ServiceDriver:     nil,
				DeploymentDriver:  &depCfgOverride,
			}
			tmp := v

			tmpSb := tmp.SkeletonBases[k]
			tmpSb.TopologyConfigDriver = &cfgOverride
			tmp.SkeletonBases[k] = tmpSb
			cd.ComponentBases[k] = tmp
		} else {
			cfgOverride := zeus_topology_config_drivers.TopologyConfigDriver{
				IngressDriver:     nil,
				StatefulSetDriver: &stsCfgOverride,
				ServiceDriver:     nil,
				DeploymentDriver:  nil,
			}
			tmp := v
			if k == "consensusClients" {
				sb := tmp.SkeletonBases["lighthouseAthena"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmp.SkeletonBases["lighthouseAthena"] = tmpSb
			} else if k == "execClients" {
				sb := tmp.SkeletonBases["gethAthena"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmp.SkeletonBases["gethAthena"] = tmpSb
			} else if k == "validatorClients" {
				sb := tmp.SkeletonBases["lighthouseAthenaValidatorClient"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmp.SkeletonBases["lighthouseAthenaValidatorClient"] = tmpSb
			}
			cd.ComponentBases[k] = tmp
		}
	}
	return cd
}
