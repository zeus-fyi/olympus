package olympus_hydra_cookbooks

import (
	"fmt"

	olympus_common_vals_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/common"
	olympus_ethereum_mev_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/mev"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/zeus/workload_config_drivers/config_overrides"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	HydraMainnet    = "hydraMainnet"
	HydraGoerli     = "hydraGoerli"
	HydraEphemery   = "hydraEphemery"
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

	consensusClientEphemeralRequestRAM      = "2Gi"
	consensusClientEphemeralRequestLimitRAM = "2Gi"

	consensusClientEphemeralRequestCPU      = "2.5"
	consensusClientEphemeralRequestLimitCPU = "2.5"

	execClientEphemeralRequestRAM      = "2Gi"
	execClientEphemeralRequestLimitRAM = "2Gi"

	execClientEphemeralRequestCPU      = "1.5"
	execClientEphemeralRequestLimitCPU = "1.5"

	consensusClientDiskName = "consensus-client-storage"
	execClientDiskName      = "exec-client-storage"

	consensusStorageDiskSizeEphemeral = "20Gi"
	execClientDiskSizeEphemeral       = "40Gi"

	gethDockerImage       = "ethereum/client-go:v1.11.5"
	lighthouseDockerImage = "sigp/lighthouse:v3.5.1"

	gethDockerImageEphemery       = "ethpandaops/geth:master"
	lighthouseDockerImageEphemery = "sigp/lighthouse:v3.5.1"

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
	switch network {
	case "mainnet":
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumMainnetProtocolNetworkID))
	case "goerli":
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumGoerliProtocolNetworkID))
	case "ephemery":
		envVar = HydraContainer.CreateEnvVarKeyValue(protocolNetworkKeyEnv, fmt.Sprintf("%d", hestia_req_types.EthereumEphemeryProtocolNetworkID))
	}
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
	combinedEnvVars = append(combinedEnvVars, envVar)

	containCfg := zeus_topology_config_drivers.ContainerDriver{}
	containCfgBeaconConsensusClient := zeus_topology_config_drivers.ContainerDriver{}
	containCfgBeaconExecClient := zeus_topology_config_drivers.ContainerDriver{}
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

	containCfgHydraClient := zeus_topology_config_drivers.ContainerDriver{
		Container: v1.Container{
			Resources: hydraRR,
			Env:       combinedEnvVars,
		},
	}

	switch network {
	case "mainnet":
		cd.CloudCtxNs.Namespace = mainnetNamespace
		cd.ClusterClassName = HydraMainnet
		cd.ComponentBases["mev"] = olympus_ethereum_mev_cookbooks.MevCbCfgMainnet()
	case "goerli":
		cd.CloudCtxNs.Namespace = goerliNamespace
		cd.ClusterClassName = HydraGoerli
		cd.ComponentBases["mev"] = olympus_ethereum_mev_cookbooks.MevCbCfgGoerli()
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

		vcContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{
			validatorClient: {Container: v1.Container{
				Name:  validatorClient,
				Image: lighthouseDockerImage,
				Env:   combinedEnvVars,
				Args:  []string{"-c", "/scripts/lighthouseGoerli" + ".sh"},
			}},
		}
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
		stsCfgOverride.ContainerDrivers[validatorClient] = vcContDriver[validatorClient]
		stsCfgOverrideSecondary.ContainerDrivers[validatorClient] = vcContDriver[validatorClient]
	case "ephemery":
		cd.CloudCtxNs.Namespace = ephemeryNamespace
		cd.ClusterClassName = HydraEphemery

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

		vcContDriver = map[string]zeus_topology_config_drivers.ContainerDriver{
			validatorClient: {Container: v1.Container{
				Name:  validatorClient,
				Image: lighthouseDockerImageEphemery,
				Env:   combinedEnvVars,
				Args:  []string{"-c", "/scripts/lighthouseEphemery" + ".sh"},
			}},
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
		stsCfgOverride.ContainerDrivers[execClient] = ecContDriver[execClient]
		stsCfgOverrideSecondary.ContainerDrivers[execClient] = ecContDriver[execClient]

		stsCfgOverride.ContainerDrivers[consensusClient] = ccContDriver[consensusClient]
		stsCfgOverrideSecondary.ContainerDrivers[consensusClient] = ccContDriver[consensusClient]

		stsCfgOverride.ContainerDrivers[validatorClient] = vcContDriver[validatorClient]
		stsCfgOverrideSecondary.ContainerDrivers[validatorClient] = vcContDriver[validatorClient]
	}
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
				sb := tmp.SkeletonBases["lighthouseAthena"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmpSb.TopologyConfigDriver.ConfigMapDriver = &ccCmDriver
				tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
				tmp.SkeletonBases["lighthouseAthena"] = tmpSb
			} else if k == "execClients" {
				tmpStsCfgOverride := stsCfgOverride
				tmpStsCfgOverride.PVCDriver = pvcEC
				sb := tmp.SkeletonBases["gethAthena"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmpSb.TopologyConfigDriver.ConfigMapDriver = &ecCmDriver
				tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
				tmp.SkeletonBases["gethAthena"] = tmpSb
			} else if k == "validatorClients" {
				tmpStsCfgOverride := stsCfgOverride
				sb := tmp.SkeletonBases["lighthouseAthenaValidatorClient"]
				tmpSb := sb
				tmpSb.TopologyConfigDriver = &cfgOverride
				tmpSb.TopologyConfigDriver.StatefulSetDriver = &tmpStsCfgOverride
				tmp.SkeletonBases["lighthouseAthenaValidatorClient"] = tmpSb
			} else if k == "validatorClientsSecondary" {
				tmpStsCfgOverride := stsCfgOverrideSecondary
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
