package olympus_beacon_cookbooks

import (
	"fmt"

	olympus_common_vals_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/common"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/zeus/workload_config_drivers/config_overrides"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	gethDockerImageEphemery       = "ethereum/client-go:v1.11.6"
	lighthouseDockerImageEphemery = "sigp/lighthouse:v4.1.0"
)

var (
	consensusClientDiskName = "consensus-client-storage"
	execClientDiskName      = "exec-client-storage"

	consensusStorageDiskSizeEphemeral = "20Gi"
	execClientDiskSizeEphemeral       = "40Gi"

	consensusClientEphemeralRequestRAM      = "1Gi"
	consensusClientEphemeralRequestLimitRAM = "1Gi"

	consensusClientEphemeralRequestCPU      = "1"
	consensusClientEphemeralRequestLimitCPU = "1"

	execClientEphemeralRequestRAM      = "1Gi"
	execClientEphemeralRequestLimitRAM = "1Gi"

	execClientEphemeralRequestCPU      = "1"
	execClientEphemeralRequestLimitCPU = "1"

	consensusClientGoerliRequestRAM      = "12Gi"
	consensusClientGoerliRequestLimitRAM = "12Gi"

	consensusClientGoerliRequestCPU      = "7"
	consensusClientGoerliRequestLimitCPU = "7"

	execClientGoerliRequestRAM      = "10Gi"
	execClientGoerliRequestLimitRAM = "10Gi"

	execClientGoerliRequestCPU      = "6"
	execClientGoerliRequestLimitCPU = "6"

	consensusStorageDiskSizeGoerli = "500Gi"
	execClientDiskSizeGoerli       = "1000Gi"
)

func ClusterConfigEnvVars(cd *zeus_cluster_config_drivers.ClusterDefinition, network string) []v1.EnvVar {
	var pvcCC *zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver
	var pvcEC *zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver
	depCfgOverride := zeus_topology_config_drivers.DeploymentDriver{}
	depCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	stsCfgOverride := zeus_topology_config_drivers.StatefulSetDriver{}
	stsCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	containCfg := zeus_topology_config_drivers.ContainerDriver{}
	containCfgBeaconConsensusClient := zeus_topology_config_drivers.ContainerDriver{}
	containCfgBeaconExecClient := zeus_topology_config_drivers.ContainerDriver{}
	envVarsChoreography := olympus_common_vals_cookbooks.GetChoreographyEnvVars()
	internalAuthEnvVars := olympus_common_vals_cookbooks.GetCommonInternalAuthEnvVars()
	combinedEnvVars := append(envVarsChoreography, internalAuthEnvVars...)
	var envVar v1.EnvVar
	switch network {
	case "mainnet":
		cdTmp := zeus_topology_config_drivers.ContainerDriver{}
		envVar = cdTmp.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
		combinedEnvVars = append(combinedEnvVars, envVar)

	case "goerli":
		rrCC := v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientGoerliRequestLimitCPU),
				"memory": resource.MustParse(consensusClientGoerliRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientGoerliRequestCPU),
				"memory": resource.MustParse(consensusClientGoerliRequestRAM),
			},
		}
		rrEC := v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(execClientGoerliRequestLimitCPU),
				"memory": resource.MustParse(execClientGoerliRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(execClientGoerliRequestCPU),
				"memory": resource.MustParse(execClientGoerliRequestRAM),
			},
		}

		envVar = containCfgBeaconConsensusClient.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumGoerliProtocolNetworkID))
		combinedEnvVars = append(combinedEnvVars, envVar)
		containCfgBeaconConsensusClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Name:      consensusClient,
				Image:     lighthouseDockerImage,
				Env:       combinedEnvVars,
				Args:      []string{"-c", "/scripts/lighthouseGoerli" + ".sh"},
				Resources: rrCC,
			},
		}
		containCfgBeaconExecClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Name:      execClient,
				Image:     gethDockerImage,
				Env:       combinedEnvVars,
				Args:      []string{"-c", "/scripts/gethGoerli" + ".sh"},
				Resources: rrEC,
			},
		}
		pvcCC = &zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver{
			PersistentVolumeClaimDrivers: map[string]v1.PersistentVolumeClaim{
				consensusClientDiskName: {
					ObjectMeta: metav1.ObjectMeta{Name: consensusClientDiskName},
					Spec: v1.PersistentVolumeClaimSpec{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{"storage": resource.MustParse(consensusStorageDiskSizeGoerli)},
					}},
				},
			}}
		pvcEC = &zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver{
			PersistentVolumeClaimDrivers: map[string]v1.PersistentVolumeClaim{
				execClientDiskName: {
					ObjectMeta: metav1.ObjectMeta{Name: execClientDiskName},
					Spec: v1.PersistentVolumeClaimSpec{Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{"storage": resource.MustParse(execClientDiskSizeGoerli)},
					}},
				},
			}}
		containCfg.Env = combinedEnvVars
	case "ephemery":
		cdTmp := zeus_topology_config_drivers.ContainerDriver{}
		envVar = cdTmp.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))
		combinedEnvVars = append(combinedEnvVars, envVar)
		rrCC := v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientEphemeralRequestLimitCPU),
				"memory": resource.MustParse(consensusClientEphemeralRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientEphemeralRequestCPU),
				"memory": resource.MustParse(consensusClientEphemeralRequestRAM),
			},
		}
		rrEC := v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(execClientEphemeralRequestLimitCPU),
				"memory": resource.MustParse(execClientEphemeralRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(execClientEphemeralRequestCPU),
				"memory": resource.MustParse(execClientEphemeralRequestRAM),
			},
		}
		containCfgBeaconConsensusClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Name:      consensusClient,
				Image:     lighthouseDockerImageEphemery,
				Env:       combinedEnvVars,
				Args:      []string{"-c", "/scripts/lighthouseEphemery" + ".sh"},
				Resources: rrCC,
			},
		}
		containCfgBeaconExecClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Name:      execClient,
				Image:     gethDockerImageEphemery,
				Env:       combinedEnvVars,
				Args:      []string{"-c", "/scripts/gethEphemery" + ".sh"},
				Resources: rrEC,
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
		containCfg.Env = combinedEnvVars
	}

	stsCfgOverride.ContainerDrivers["init-snapshots"] = containCfg
	stsCfgOverride.ContainerDrivers[athena] = containCfg
	stsCfgOverride.ContainerDrivers[execClient] = containCfgBeaconExecClient
	stsCfgOverride.ContainerDrivers[consensusClient] = containCfgBeaconConsensusClient

	if cd == nil {
		return combinedEnvVars
	}

	for k, v := range cd.ComponentBases {
		cfgOverride := zeus_topology_config_drivers.TopologyConfigDriver{
			IngressDriver:     nil,
			StatefulSetDriver: &stsCfgOverride,
			ServiceDriver:     nil,
			DeploymentDriver:  nil,
		}
		tmp := v
		if k == "consensusClients" {
			tmpStsCfgOverride := stsCfgOverride
			tmpStsCfgOverride.PVCDriver = pvcCC
			sb := tmp.SkeletonBases["lighthouseAthena"]
			tmpSb := sb
			tmpSb.TopologyConfigDriver = &cfgOverride
			tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
			tmp.SkeletonBases["lighthouseAthena"] = tmpSb
		} else if k == "execClients" {
			tmpStsCfgOverride := stsCfgOverride
			tmpStsCfgOverride.PVCDriver = pvcEC
			sb := tmp.SkeletonBases["gethAthena"]
			tmpSb := sb
			tmpSb.TopologyConfigDriver = &cfgOverride
			tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
			tmp.SkeletonBases["gethAthena"] = tmpSb
		}
		cd.ComponentBases[k] = tmp
	}

	return combinedEnvVars
}
