package zeus_templates

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_resp_types/topology_workloads"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func PreviewTemplateGeneration(ctx context.Context, cluster Cluster) zeus_cluster_config_drivers.ClusterDefinition {
	templateClusterDefinition := zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: cluster.ClusterName,
		ComponentBases:   make(map[string]zeus_cluster_config_drivers.ComponentBaseDefinition),
	}
	fmt.Println(templateClusterDefinition)
	for cbName, componentBase := range cluster.ComponentBases {
		cbDef := zeus_cluster_config_drivers.ComponentBaseDefinition{
			SkeletonBases: make(map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition),
		}
		for sbName, skeletonBase := range componentBase {
			sbDef := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
				SkeletonBaseChart: zeus_req_types.TopologyCreateRequest{},
				SkeletonBaseNameChartPath: filepaths.Path{
					PackageName: sbName,
					DirIn:       "./",
					DirOut:      "./",
					FnIn:        sbName,
					FnOut:       sbName,
					FilterFiles: &strings_filter.FilterOpts{},
				},
				Workload:             topology_workloads.TopologyBaseInfraWorkload{},
				TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{},
			}
			if skeletonBase.AddStatefulSet {
				stsDriver, _ := BuildStatefulSetDriver(ctx, sbName, skeletonBase.Containers, skeletonBase.StatefulSet)
				sbDef.TopologyConfigDriver.StatefulSetDriver = &stsDriver
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "deployment")
			} else if skeletonBase.AddDeployment {
				depDriver, _ := BuildDeploymentDriver(ctx, sbName, skeletonBase.Containers, skeletonBase.Deployment)
				sbDef.TopologyConfigDriver.DeploymentDriver = &depDriver
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "statefulset")
			}
			if skeletonBase.AddIngress {
				ingDriver, _ := BuildIngressDriver(ctx, sbName, cluster.IngressSettings, cluster.IngressPaths)
				sbDef.TopologyConfigDriver.IngressDriver = &ingDriver
			} else {
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "ingress")
			}
			if skeletonBase.AddService {
				svcDriver, _ := BuildServiceDriver(ctx, sbName, skeletonBase.Containers)
				sbDef.TopologyConfigDriver.ServiceDriver = &svcDriver
			} else {
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "service")
			}
			if skeletonBase.AddConfigMap {
				cmDriver, _ := BuildConfigMapDriver(ctx, sbName, skeletonBase.ConfigMap)
				sbDef.TopologyConfigDriver.ConfigMapDriver = &cmDriver
			} else {
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "configmap")
			}
			cbDef.SkeletonBases[sbName] = sbDef
		}
		templateClusterDefinition.ComponentBases[cbName] = cbDef
	}
	return templateClusterDefinition
}

func BuildStatefulSetDriver(ctx context.Context, sbName string, containers Containers, sts StatefulSet) (zeus_topology_config_drivers.StatefulSetDriver, error) {
	rc := int32(sts.ReplicaCount)
	stsDriver := zeus_topology_config_drivers.StatefulSetDriver{
		ReplicaCount:     &rc,
		ContainerDrivers: make(map[string]zeus_topology_config_drivers.ContainerDriver),
	}
	for containerName, container := range containers {
		contDriver, err := BuildContainerDriver(ctx, sbName, container)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to build container driver")
			return zeus_topology_config_drivers.StatefulSetDriver{}, err
		}
		stsDriver.ContainerDrivers[containerName] = zeus_topology_config_drivers.ContainerDriver{
			IsAppendContainer: true,
			IsInitContainer:   container.IsInitContainer,
			Container:         contDriver,
			AppendEnvVars:     nil,
		}
	}
	pvcCfg := zeus_topology_config_drivers.PersistentVolumeClaimsConfigDriver{
		AppendPVC:                    make(map[string]bool),
		PersistentVolumeClaimDrivers: make(map[string]v1.PersistentVolumeClaim),
	}
	for _, pvcTemplate := range sts.PVCTemplates {
		storageReq := v1.ResourceList{"storage": resource.MustParse(pvcTemplate.StorageSizeRequest)}
		pvc := v1.PersistentVolumeClaim{
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{v1.PersistentVolumeAccessMode(pvcTemplate.AccessMode)},
				Resources: v1.ResourceRequirements{
					Requests: storageReq,
				},
				VolumeName: pvcTemplate.Name,
			},
		}
		pvcCfg.AppendPVC[pvcTemplate.Name] = true
		pvcCfg.PersistentVolumeClaimDrivers[pvcTemplate.Name] = pvc
	}
	stsDriver.PVCDriver = &pvcCfg
	return stsDriver, nil
}

func BuildDeploymentDriver(ctx context.Context, sbName string, containers Containers, dep Deployment) (zeus_topology_config_drivers.DeploymentDriver, error) {
	rc := int32(dep.ReplicaCount)
	depDriver := zeus_topology_config_drivers.DeploymentDriver{
		ReplicaCount:     &rc,
		ContainerDrivers: make(map[string]zeus_topology_config_drivers.ContainerDriver),
	}
	for containerName, container := range containers {
		contDriver, err := BuildContainerDriver(ctx, sbName, container)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to build container driver")
			return zeus_topology_config_drivers.DeploymentDriver{}, err
		}
		depDriver.ContainerDrivers[containerName] = zeus_topology_config_drivers.ContainerDriver{
			IsAppendContainer: true,
			IsInitContainer:   container.IsInitContainer,
			Container:         contDriver,
			AppendEnvVars:     nil,
		}
	}
	return depDriver, nil
}

func BuildServiceDriver(ctx context.Context, sbName string, containers Containers) (zeus_topology_config_drivers.ServiceDriver, error) {
	svcDriver := zeus_topology_config_drivers.ServiceDriver{
		Service: v1.Service{
			Spec: v1.ServiceSpec{
				Ports: []v1.ServicePort{},
			},
		},
	}
	var sps []v1.ServicePort
	for _, c := range containers {
		for _, p := range c.DockerImage.Ports {
			numberInt64, err := strconv.ParseInt(p.Number, 10, 32)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to parse port number")
				return svcDriver, err
			}
			sps = append(sps, v1.ServicePort{
				Name:       p.Name,
				Port:       int32(numberInt64),
				Protocol:   v1.Protocol(p.Protocol),
				TargetPort: intstr.IntOrString{Type: intstr.String, StrVal: p.Name},
			})
			if p.IngressEnabledPort {
				svcDriver.AddNginxTargetPort("http", p.Name)
			}
		}
	}
	svcDriver.Service.Spec.Ports = sps
	return svcDriver, nil
}

func BuildIngressDriver(ctx context.Context, sbName string, ing Ingress, ip IngressPaths) (zeus_topology_config_drivers.IngressDriver, error) {
	var httpPaths []v1networking.HTTPIngressPath
	for _, pa := range ip {
		pt := v1networking.PathType(pa.PathType)
		appendPath := v1networking.HTTPIngressPath{
			Path:     pa.Path,
			PathType: &pt,
			Backend: v1networking.IngressBackend{
				Service: &v1networking.IngressServiceBackend{
					Name: "http", // TODO rename
					Port: v1networking.ServiceBackendPort{
						Number: int32(80),
					},
				},
			},
		}
		httpPaths = append(httpPaths, appendPath)
	}
	ingressRuleValue := v1networking.IngressRuleValue{HTTP: &v1networking.HTTPIngressRuleValue{Paths: httpPaths}}
	ingDriver := zeus_topology_config_drivers.IngressDriver{
		Ingress: v1networking.Ingress{
			Spec: v1networking.IngressSpec{
				TLS: []v1networking.IngressTLS{{
					Hosts:      []string{ing.Host},
					SecretName: "tls-secret", // TODO rename
				}},
				Rules: []v1networking.IngressRule{{
					Host:             ing.Host,
					IngressRuleValue: ingressRuleValue,
				}},
			},
		},
		Host:         ing.Host,
		NginxAuthURL: ing.AuthServerURL,
	}
	return ingDriver, nil
}

func BuildConfigMapDriver(ctx context.Context, sbName string, configMap ConfigMap) (zeus_topology_config_drivers.ConfigMapDriver, error) {
	cmDriver := zeus_topology_config_drivers.ConfigMapDriver{
		ConfigMap: v1.ConfigMap{
			Data: make(map[string]string),
		},
	}
	for key, value := range configMap {
		cmDriver.ConfigMap.Data[key] = value
	}
	return cmDriver, nil
}

func LabelBuilder(ctx context.Context) {
	// TODO
}

func BuildContainerDriver(ctx context.Context, sbName string, container Container) (v1.Container, error) {
	pp := "IfNotPresent"
	if len(container.ImagePullPolicy) <= 0 {
		pp = container.ImagePullPolicy
	}
	c := v1.Container{
		Name:    sbName,
		Image:   container.DockerImage.ImageName,
		Command: strings.Split(container.DockerImage.Cmd, ","),
		Args:    strings.Split(container.DockerImage.Args, ","),
		Ports:   []v1.ContainerPort{},
		EnvFrom: nil,
		Env:     nil,
		Resources: v1.ResourceRequirements{
			Limits:   make(map[v1.ResourceName]resource.Quantity),
			Requests: make(map[v1.ResourceName]resource.Quantity),
		},
		VolumeMounts:    []v1.VolumeMount{},
		ImagePullPolicy: v1.PullPolicy(pp),
	}

	for _, p := range container.DockerImage.Ports {
		// Use strconv.ParseInt to convert the string to int64
		numberInt64, err := strconv.ParseInt(p.Number, 10, 32)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to parse port number")
			return c, err
		}
		c.Ports = append(c.Ports, v1.ContainerPort{
			Name:          p.Name,
			ContainerPort: int32(numberInt64),
			Protocol:      v1.Protocol(p.Protocol),
		})
	}

	for _, v := range container.DockerImage.VolumeMounts {
		c.VolumeMounts = append(c.VolumeMounts, v1.VolumeMount{
			Name:      v.Name,
			MountPath: v.MountPath,
		})
	}
	if len(container.DockerImage.ResourceRequirements.CPU) > 0 {
		c.Resources.Requests["cpu"] = resource.MustParse(container.DockerImage.ResourceRequirements.CPU)
		c.Resources.Limits["cpu"] = resource.MustParse(container.DockerImage.ResourceRequirements.CPU)
	}
	if len(container.DockerImage.ResourceRequirements.Memory) > 0 {
		c.Resources.Requests["memory"] = resource.MustParse(container.DockerImage.ResourceRequirements.Memory)
		c.Resources.Limits["memory"] = resource.MustParse(container.DockerImage.ResourceRequirements.Memory)
	}
	return c, nil
}
