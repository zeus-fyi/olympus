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
	consensusClient = "zeus-consensus-client"
	execClient      = "zeus-exec-client"
	validatorClient = "zeus-hydra-validators"

	protocolNetworkKeyEnv = "PROTOCOL_NETWORK_ID"
	ephemeryNamespace     = "ephemeral-staking"
	goerliNamespace       = "goerli-staking"
	mainnetNamespace      = "mainnet-staking"

	hydraClientEphemeralRequestRAM      = "500Mi"
	hydraClientEphemeralRequestLimitRAM = "500Mi"

	hydraClientEphemeralRequestCPU      = "2.5"
	hydraClientEphemeralRequestLimitCPU = "2.5"

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

	consensusStorageDiskSizeEphemeral = "20Gi"
	execClientDiskSizeEphemeral       = "40Gi"

	gethDockerImage       = "ethereum/client-go:v1.11.4"
	lighthouseDockerImage = "sigp/lighthouse:v3.3.0-modern"

	cmExecClient      = "cm-exec-client"
	cmConsensusClient = "cm-consensus-client"
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
	var hydraRR v1.ResourceRequirements

	hydraRR = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu":    resource.MustParse(hydraClientEphemeralRequestCPU),
			"memory": resource.MustParse(hydraClientEphemeralRequestLimitRAM),
		},
		Requests: v1.ResourceList{
			"cpu":    resource.MustParse(hydraClientEphemeralRequestLimitCPU),
			"memory": resource.MustParse(hydraClientEphemeralRequestRAM),
		},
	}

	var pvcCC *zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver
	var pvcEC *zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver
	var ccContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{}
	var ecContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{}
	var vcContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{}
	var ccCmDriver = zeus_topology_config_drivers.ConfigMapDriver{}
	var ecCmDriver = zeus_topology_config_drivers.ConfigMapDriver{}

	depCfgOverride := zeus_topology_config_drivers.DeploymentDriver{}
	depCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	stsCfgOverride := zeus_topology_config_drivers.StatefulSetDriver{}
	stsCfgOverride.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)
	stsCfgOverrideSecondary := zeus_topology_config_drivers.StatefulSetDriver{}
	stsCfgOverrideSecondary.ContainerDrivers = make(map[string]zeus_topology_config_drivers.ContainerDriver)

	envVarsChoreography := olympus_common_vals_cookbooks.GetChoreographyEnvVars()
	internalAuthEnvVars := olympus_common_vals_cookbooks.GetCommonInternalAuthEnvVars()
	combinedEnvVars := append(envVarsChoreography, internalAuthEnvVars...)

	containCfg := zeus_topology_config_drivers.ContainerDriver{}
	containCfgBeaconConsensusClient := zeus_topology_config_drivers.ContainerDriver{}
	containCfgBeaconExecClient := zeus_topology_config_drivers.ContainerDriver{}

	switch network {
	case "mainnet":
		cd.CloudCtxNs.Namespace = mainnetNamespace
		cd.ClusterClassName = "hydraMainnet"
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
	case "goerli":
		cd.CloudCtxNs.Namespace = goerliNamespace
		cd.ClusterClassName = "hydraGoerli"
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumGoerliProtocolNetworkID))
		combinedEnvVars = append(combinedEnvVars, envVar)

		rrCC = v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientGoerliRequestLimitCPU),
				"memory": resource.MustParse(consensusClientGoerliRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(consensusClientGoerliRequestCPU),
				"memory": resource.MustParse(consensusClientGoerliRequestRAM),
			},
		}
		rrEC = v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(execClientGoerliRequestLimitCPU),
				"memory": resource.MustParse(execClientGoerliRequestLimitRAM),
			},
			Requests: v1.ResourceList{
				"cpu":    resource.MustParse(execClientGoerliRequestCPU),
				"memory": resource.MustParse(execClientGoerliRequestRAM),
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
		ccContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{
			consensusClient: {Container: v1.Container{
				Name:  consensusClient,
				Image: lighthouseDockerImage,
				Env:   combinedEnvVars,
				Args:  []string{"-c", "/scripts/lighthouseGoerli" + ".sh"},
			}},
		}
		ecContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{
			execClient: {Container: v1.Container{
				Name:  execClient,
				Image: gethDockerImage,
				Env:   combinedEnvVars,
				Args:  []string{"-c", "/scripts/gethGoerli" + ".sh"},
			}},
		}

		vcContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{
			validatorClient: {Container: v1.Container{
				Name:  validatorClient,
				Image: lighthouseDockerImage,
				Env:   combinedEnvVars,
				Args:  []string{"-c", "/scripts/lighthouseGoerli" + ".sh"},
			}},
		}
		ecCmDriver = zeus_topology_config_drivers.ConfigMapDriver{
			ConfigMap: v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: cmExecClient},
			},
		}
		containCfgBeaconConsensusClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Resources: rrCC,
			},
		}

		containCfgBeaconExecClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Resources: rrEC,
			},
		}
	case "ephemery":
		cd.CloudCtxNs.Namespace = ephemeryNamespace
		cd.ClusterClassName = "hydraEphemery"
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))
		combinedEnvVars = append(combinedEnvVars, envVar)

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
		containCfgBeaconConsensusClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Resources: rrCC,
			},
		}

		containCfgBeaconExecClient = zeus_topology_config_drivers.ContainerDriver{
			Container: v1.Container{
				Resources: rrEC,
			},
		}
	}

	containCfgHydraClient := zeus_topology_config_drivers.ContainerDriver{
		Container: v1.Container{
			Resources: hydraRR,
			Env:       combinedEnvVars,
		},
	}

	containCfg.Env = combinedEnvVars
	containCfgSecondary := containCfg
	rcSecondary := v1.EnvVar{
		Name:  "REPLICA_COUNT",
		Value: "1",
	}
	containCfgSecondary.AppendEnvVars = []v1.EnvVar{rcSecondary}

	rcPrimary := v1.EnvVar{
		Name:  "REPLICA_COUNT",
		Value: "0",
	}
	containCfg.AppendEnvVars = []v1.EnvVar{rcPrimary}

	// deployments
	depCfgOverride.ContainerDrivers["hydra"] = containCfgHydraClient
	depCfgOverride.ContainerDrivers["zeus-hydra-choreography"] = containCfg
	depCfgOverride.ContainerDrivers["athena"] = containCfg

	// statefulsets
	stsCfgOverride.ContainerDrivers["athena"] = containCfg
	stsCfgOverride.ContainerDrivers["zeus-consensus-client"] = containCfgBeaconConsensusClient
	stsCfgOverride.ContainerDrivers["zeus-exec-client"] = containCfgBeaconExecClient
	stsCfgOverride.ContainerDrivers["init-validators"] = containCfg
	stsCfgOverride.ContainerDrivers["init-snapshots"] = containCfg

	stsCfgOverrideSecondary.ContainerDrivers["athena"] = containCfgSecondary
	stsCfgOverrideSecondary.ContainerDrivers["zeus-consensus-client"] = containCfgBeaconConsensusClient
	stsCfgOverrideSecondary.ContainerDrivers["zeus-exec-client"] = containCfgBeaconExecClient
	stsCfgOverrideSecondary.ContainerDrivers["init-validators"] = containCfgSecondary
	stsCfgOverrideSecondary.ContainerDrivers["init-snapshots"] = containCfgSecondary

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
				tmpStsCfgOverride := stsCfgOverride
				tmpStsCfgOverride.PVCDriver = pvcCC
				tmpStsCfgOverride.ContainerDrivers = ccContDriver
				sb := tmp.SkeletonBases["lighthouseAthena"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmpSb.TopologyConfigDriver.ConfigMapDriver = &ccCmDriver
				tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
				tmp.SkeletonBases["lighthouseAthena"] = tmpSb
			} else if k == "execClients" {
				tmpStsCfgOverride := stsCfgOverride
				tmpStsCfgOverride.PVCDriver = pvcEC
				tmpStsCfgOverride.ContainerDrivers = ecContDriver
				sb := tmp.SkeletonBases["gethAthena"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmpSb.TopologyConfigDriver.ConfigMapDriver = &ecCmDriver
				tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
				tmp.SkeletonBases["gethAthena"] = tmpSb
			} else if k == "validatorClients" {
				tmpStsCfgOverride := stsCfgOverride
				tmpStsCfgOverride.ContainerDrivers = vcContDriver
				sb := tmp.SkeletonBases["lighthouseAthenaValidatorClient"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
				tmp.SkeletonBases["lighthouseAthenaValidatorClient"] = tmpSb
			} else if k == "validatorClientsSecondary" {
				tmpStsCfgOverride := stsCfgOverrideSecondary
				tmpStsCfgOverride.ContainerDrivers = vcContDriver
				sb := tmp.SkeletonBases["lighthouseAthenaValidatorClientSecondary"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
				tmp.SkeletonBases["lighthouseAthenaValidatorClientSecondary"] = tmpSb
			}
			cd.ComponentBases[k] = tmp
		}
	}
	return cd
}
