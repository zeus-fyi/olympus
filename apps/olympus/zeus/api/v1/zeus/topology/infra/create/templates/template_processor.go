package zeus_templates

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_resp_types/topology_workloads"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ClusterPreviewWorkloads struct {
	ClusterName                   string                                                             `json:"clusterName"`
	ComponentBasesToSkeletonBases map[string]map[string]topology_workloads.TopologyBaseInfraWorkload `json:"componentBases"`
}

func GenerateSkeletonBaseChartsPreview(ctx context.Context, cluster Cluster) (ClusterPreviewWorkloads, error) {
	pcg := ClusterPreviewWorkloads{
		ClusterName:                   cluster.ClusterName,
		ComponentBasesToSkeletonBases: make(map[string]map[string]topology_workloads.TopologyBaseInfraWorkload),
	}
	cd := PreviewTemplateGeneration(ctx, cluster)
	cd.UseEmbeddedWorkload = true
	cd.DisablePrint = true
	_, err := cd.GenerateSkeletonBaseCharts()
	if err != nil {
		log.Ctx(ctx).Err(err)
		return pcg, err
	}
	for cbName, componentBase := range cd.ComponentBases {
		pcg.ComponentBasesToSkeletonBases[cbName] = make(map[string]topology_workloads.TopologyBaseInfraWorkload)
		for sbName, skeletonBase := range componentBase.SkeletonBases {
			pcg.ComponentBasesToSkeletonBases[cbName][sbName] = skeletonBase.Workload
		}
	}
	return pcg, nil
}

func GenerateClusterFromUI(ctx context.Context, cluster Cluster) (zeus_cluster_config_drivers.GeneratedClusterCreationRequests, error) {
	cd := PreviewTemplateGeneration(ctx, cluster)
	_, err := cd.GenerateSkeletonBaseCharts()
	if err != nil {
		log.Ctx(ctx).Err(err)
		return zeus_cluster_config_drivers.GeneratedClusterCreationRequests{}, err
	}
	gcd := cd.BuildClusterDefinitions()
	if err != nil {
		log.Ctx(ctx).Err(err)
		return gcd, err
	}
	return gcd, nil
}

func PreviewTemplateGeneration(ctx context.Context, cluster Cluster) zeus_cluster_config_drivers.ClusterDefinition {
	templateClusterDefinition := zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: cluster.ClusterName,
		ComponentBases:   make(map[string]zeus_cluster_config_drivers.ComponentBaseDefinition),
	}
	for cbName, componentBase := range cluster.ComponentBases {
		cbDef := zeus_cluster_config_drivers.ComponentBaseDefinition{
			SkeletonBases: make(map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition),
		}
		for sbName, skeletonBase := range componentBase {
			sbDef := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
				SkeletonBaseChart:    zeus_req_types.TopologyCreateRequest{},
				Workload:             topology_workloads.TopologyBaseInfraWorkload{},
				TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{},
			}
			if skeletonBase.AddStatefulSet {
				sbDef.Workload.StatefulSet = GetStatefulSetTemplate(ctx, cbName)
				stsDriver, err := BuildStatefulSetDriver(ctx, sbName, skeletonBase.Containers, skeletonBase.StatefulSet)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("error building statefulset driver")
				}
				sbDef.TopologyConfigDriver.StatefulSetDriver = &stsDriver
			} else if skeletonBase.AddDeployment {
				sbDef.Workload.Deployment = GetDeploymentTemplate(ctx, cbName)
				depDriver, err := BuildDeploymentDriver(ctx, sbName, skeletonBase.Containers, skeletonBase.Deployment)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("error building deployment driver")
				}
				sbDef.TopologyConfigDriver.DeploymentDriver = &depDriver
			}
			if skeletonBase.AddIngress {
				sbDef.Workload.Ingress = GetIngressTemplate(ctx, cbName)
				ingDriver, err := BuildIngressDriver(ctx, cbName, skeletonBase.Containers, cluster.IngressSettings, cluster.IngressPaths)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("error building ingress driver")
				}
				sbDef.TopologyConfigDriver.IngressDriver = &ingDriver
			}
			if skeletonBase.AddService {
				sbDef.Workload.Service = GetServiceTemplate(ctx, cbName)
				svcDriver, err := BuildServiceDriver(ctx, skeletonBase.Containers)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("error building service driver")
				}
				sbDef.TopologyConfigDriver.ServiceDriver = &svcDriver
			}
			if skeletonBase.AddConfigMap {
				sbDef.Workload.ConfigMap = GetConfigMapTemplate(ctx, cbName)
				cmDriver, err := BuildConfigMapDriver(ctx, skeletonBase.ConfigMap)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("error building configmap driver")
				}
				sbDef.TopologyConfigDriver.ConfigMapDriver = &cmDriver
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

func BuildServiceDriver(ctx context.Context, containers Containers) (zeus_topology_config_drivers.ServiceDriver, error) {
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

func BuildIngressDriver(ctx context.Context, cbName string, containers Containers, ing Ingress, ip IngressPaths) (zeus_topology_config_drivers.IngressDriver, error) {
	portName := ""
	uid := uuid.New()
	ing.Host = GetIngressHostName(ctx, uid.String())
	for _, container := range containers {
		for _, p := range container.DockerImage.Ports {
			if p.IngressEnabledPort {
				portName = p.Name
			}
		}
	}
	var httpPaths []v1networking.HTTPIngressPath
	for _, pa := range ip {
		pt := v1networking.PathType(pa.PathType)
		appendPath := v1networking.HTTPIngressPath{
			Path:     pa.Path,
			PathType: &pt,
			Backend: v1networking.IngressBackend{
				Service: &v1networking.IngressServiceBackend{
					Name: GetServiceName(ctx, cbName),
					Port: v1networking.ServiceBackendPort{
						Number: int32(80),
						Name:   portName,
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
					SecretName: GetIngressSecretName(ctx, uid.String()),
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

func BuildConfigMapDriver(ctx context.Context, configMap ConfigMap) (zeus_topology_config_drivers.ConfigMapDriver, error) {
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
