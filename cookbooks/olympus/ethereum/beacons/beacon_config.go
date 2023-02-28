package olympus_beacon_cookbooks

import (
	"fmt"

	olympus_common_vals_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/common"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
)

var (
	consensusClientStsEphemeralCfg = zeus_topology_config_drivers.StatefulSetDriver{
		ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
			"zeus-consensus-client": {
				Container: v1.Container{
					Name:  "zeus-consensus-client",
					Image: "sigp/lighthouse:capella",
				},
			},
		}}
	consensusClientStsEphemeralCfgDriver = zeus_topology_config_drivers.TopologyConfigDriver{
		StatefulSetDriver: &consensusClientStsEphemeralCfg,
	}
	execClientStsEphemeralCfg = zeus_topology_config_drivers.StatefulSetDriver{
		ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
			"zeus-exec-client": {
				Container: v1.Container{
					Name:  "zeus-exec-client",
					Image: "ethpandaops/geth:master",
				},
			},
		}}
	execClientStsEphemeralCfgDriver = zeus_topology_config_drivers.TopologyConfigDriver{
		StatefulSetDriver: &execClientStsEphemeralCfg,
	}
)

func ClusterConfigEnvVars(cd *zeus_cluster_config_drivers.ClusterDefinition, network string) []v1.EnvVar {
	var envVar v1.EnvVar
	switch network {
	case "mainnet":
		cdTmp := zeus_topology_config_drivers.ContainerDriver{}
		envVar = cdTmp.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
	case "ephemery":
		cdTmp := zeus_topology_config_drivers.ContainerDriver{}
		envVar = cdTmp.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))
	}

	depCfgOverride := zeus_topology_config_drivers.DeploymentDriver{}
	depCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	stsCfgOverride := zeus_topology_config_drivers.StatefulSetDriver{}
	stsCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	envVarsChoreography := olympus_common_vals_cookbooks.GetChoreographyEnvVars()
	internalAuthEnvVars := olympus_common_vals_cookbooks.GetCommonInternalAuthEnvVars()
	combinedEnvVars := append(envVarsChoreography, internalAuthEnvVars...)
	combinedEnvVars = append(combinedEnvVars, envVar)

	if cd == nil {
		return combinedEnvVars
	}
	for k, v := range cd.ComponentBases {
		if k == "hydra" || k == "hydraChoreography" {
			containCfg := zeus_topology_config_drivers.ContainerDriver{}
			containCfg.Env = combinedEnvVars

			// deployments
			depCfgOverride.ContainerDrivers["zeus-hydra-choreography"] = containCfg
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
		}
	}
	return combinedEnvVars
}
