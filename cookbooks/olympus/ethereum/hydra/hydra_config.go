package olympus_hydra_cookbooks

import (
	"fmt"

	olympus_common_vals_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/common"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	protocolNetworkKeyEnv = "PROTOCOL_NETWORK_ID"
	ephemeryNamespace     = "ephemeral-staking"
	mainnetNamespace      = "mainnet-staking"

	consensusClientEphemeralRequestRAM      = "1Gi"
	consensusClientEphemeralRequestLimitRAM = "1Gi"

	consensusClientEphemeralRequestCPU      = "1"
	consensusClientEphemeralRequestLimitCPU = "1"

	execClientEphemeralRequestRAM      = "1Gi"
	execClientEphemeralRequestLimitRAM = "1Gi"

	execClientEphemeralRequestCPU      = "1"
	execClientEphemeralRequestLimitCPU = "1"

	consensusClientDiskName = "consensus-client-storage"
	execClientDiskName      = "exec-client-storage"

	consensusStorageDiskSizeEphemeral = "4Gi"
	execClientDiskSizeEphemeral       = "12Gi"
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
	var rrCC v1.ResourceRequirements
	var rrEC v1.ResourceRequirements
	var pvcCC *zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver
	var pvcEC *zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver
	switch network {
	case "mainnet":
		cd.CloudCtxNs.Namespace = mainnetNamespace
		cd.ClusterClassName = "hydraMainnet"
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
	case "ephemery":
		cd.CloudCtxNs.Namespace = ephemeryNamespace
		cd.ClusterClassName = "hydraEphemery"
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))

		rrCC = v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientEphemeralRequestLimitCPU),
				"memory": resource.MustParse(consensusClientEphemeralRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientEphemeralRequestCPU),
				"memory": resource.MustParse(consensusClientEphemeralRequestRAM),
			},
		}
		rrEC = v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(execClientEphemeralRequestLimitCPU),
				"memory": resource.MustParse(execClientEphemeralRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(execClientEphemeralRequestCPU),
				"memory": resource.MustParse(execClientEphemeralRequestRAM),
			},
		}
		pvcCC = &zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver{
			PersistentVolumeClaimDrivers: map[string]v1.PersistentVolumeClaim{
				consensusClientDiskName: {
					ObjectMeta: metav1.ObjectMeta{Name: consensusClientDiskName},
					Spec: v1.PersistentVolumeClaimSpec{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{"storage": resource.MustParse(consensusStorageDiskSizeEphemeral)},
					}},
				},
			}}
		pvcEC = &zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver{
			PersistentVolumeClaimDrivers: map[string]v1.PersistentVolumeClaim{
				execClientDiskName: {
					ObjectMeta: metav1.ObjectMeta{Name: execClientDiskName},
					Spec: v1.PersistentVolumeClaimSpec{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{"storage": resource.MustParse(execClientDiskSizeEphemeral)},
					}},
				},
			}}
	}

	depCfgOverride := zeus_topology_config_drivers.DeploymentDriver{}
	depCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	stsCfgOverride := zeus_topology_config_drivers.StatefulSetDriver{}
	stsCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)

	envVarsChoreography := olympus_common_vals_cookbooks.GetChoreographyEnvVars()
	internalAuthEnvVars := olympus_common_vals_cookbooks.GetCommonInternalAuthEnvVars()
	combinedEnvVars := append(envVarsChoreography, internalAuthEnvVars...)
	combinedEnvVars = append(combinedEnvVars, envVar)

	containCfg := zeus_topology_config_drivers.ContainerDriver{}

	containCfgBeaconConsensusClient := zeus_topology_config_drivers.ContainerDriver{
		Container: v1.Container{
			Resources: rrCC,
		},
	}

	containCfgBeaconExecClient := zeus_topology_config_drivers.ContainerDriver{
		Container: v1.Container{
			Resources: rrEC,
		},
	}

	containCfg.Env = combinedEnvVars

	// deployments
	depCfgOverride.ContainerDrivers["hydra"] = containCfg
	depCfgOverride.ContainerDrivers["zeus-hydra-choreography"] = containCfg
	depCfgOverride.ContainerDrivers["athena"] = containCfg

	// statefulsets
	stsCfgOverride.ContainerDrivers["athena"] = containCfg
	stsCfgOverride.ContainerDrivers["zeus-consensus-client"] = containCfgBeaconConsensusClient
	stsCfgOverride.ContainerDrivers["zeus-exec-client"] = containCfgBeaconExecClient
	stsCfgOverride.ContainerDrivers["init-validators"] = containCfg
	stsCfgOverride.ContainerDrivers["init-snapshots"] = containCfg

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
				stsCfgOverride.PVCDriver = pvcCC
				sb := tmp.SkeletonBases["lighthouseAthena"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmp.SkeletonBases["lighthouseAthena"] = tmpSb
			} else if k == "execClients" {
				stsCfgOverride.PVCDriver = pvcEC
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
