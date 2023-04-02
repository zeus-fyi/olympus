package zeus_templates

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_resp_types/topology_workloads"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/cluster_config_drivers"
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PreviewTemplateGeneration(ctx context.Context, cluster create_infra.Cluster) {
	templateClusterDefinition := zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: cluster.ClusterName,
		ComponentBases:   make(map[string]zeus_cluster_config_drivers.ComponentBaseDefinition),
	}
	fmt.Println(templateClusterDefinition)
	for cbName, componentBase := range cluster.ComponentBases {
		fmt.Println(cbName)
		fmt.Println(componentBase)
		cbDef := zeus_cluster_config_drivers.ComponentBaseDefinition{}

		for sbName, skeletonBase := range componentBase {
			fmt.Println(sbName)
			fmt.Println(skeletonBase)
			sbDef := zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
				SkeletonBaseChart: zeus_req_types.TopologyCreateRequest{},
				SkeletonBaseNameChartPath: filepaths.Path{
					PackageName: sbName,
					DirIn:       "./templates",
					DirOut:      "",
					FnIn:        sbName,
					FnOut:       sbName,
					FilterFiles: &strings_filter.FilterOpts{},
				},
				Workload:             topology_workloads.TopologyBaseInfraWorkload{},
				TopologyConfigDriver: &zeus_topology_config_drivers.TopologyConfigDriver{},
			}

			if skeletonBase.AddStatefulSet {
				stsDriver, _ := BuildStatefulSetDriver(ctx, sbName, skeletonBase.Containers, skeletonBase.StatefulSet)
				fmt.Println(stsDriver)
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "deployment")
			} else if skeletonBase.AddDeployment {
				depDriver, _ := BuildDeploymentDriver(ctx, sbName, skeletonBase.Containers, skeletonBase.Deployment)
				fmt.Println(depDriver)
				sbDef.TopologyConfigDriver.DeploymentDriver = &depDriver
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "statefulset")
			}
			if skeletonBase.AddIngress {
				ingDriver, _ := BuildIngressDriver(ctx, sbName, cluster.IngressSettings, cluster.IngressPaths)
				fmt.Println(ingDriver)
				sbDef.TopologyConfigDriver.IngressDriver = &ingDriver
			} else {
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "ingress")
			}
			if skeletonBase.AddService {
				svcDriver, _ := BuildServiceDriver(ctx, sbName)
				fmt.Println(svcDriver)
				sbDef.TopologyConfigDriver.ServiceDriver = &svcDriver
			} else {
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "service")
			}
			if skeletonBase.AddConfigMap {
				cmDriver, _ := BuildConfigMapDriver(ctx, sbName, skeletonBase.ConfigMap)
				fmt.Println(cmDriver)
				sbDef.TopologyConfigDriver.ConfigMapDriver = &cmDriver
			} else {
				sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese = append(sbDef.SkeletonBaseNameChartPath.FilterFiles.DoesNotStartWithThese, "configmap")
			}
			cbDef.SkeletonBases[sbName] = sbDef
		}
		templateClusterDefinition.ComponentBases[cbName] = cbDef
	}
	fmt.Println(templateClusterDefinition)
}

func BuildStatefulSetDriver(ctx context.Context, sbName string, containers create_infra.Containers, sts create_infra.StatefulSet) (zeus_topology_config_drivers.StatefulSetDriver, error) {
	stsDriver := zeus_topology_config_drivers.StatefulSetDriver{}
	return stsDriver, nil
}

func BuildDeploymentDriver(ctx context.Context, sbName string, containers create_infra.Containers, dep create_infra.Deployment) (zeus_topology_config_drivers.DeploymentDriver, error) {
	rc := int32(dep.ReplicaCount)
	depDriver := zeus_topology_config_drivers.DeploymentDriver{
		ReplicaCount:     &rc,
		ContainerDrivers: make(map[string]zeus_topology_config_drivers.ContainerDriver),
	}
	for containerName, container := range containers {
		fmt.Println(containerName)
		fmt.Println(container)
		contDriver, _ := BuildContainerDriver(ctx, sbName, container)
		fmt.Println(contDriver)
	}
	return depDriver, nil
}

func BuildServiceDriver(ctx context.Context, sbName string) (zeus_topology_config_drivers.ServiceDriver, error) {
	svcDriver := zeus_topology_config_drivers.ServiceDriver{}
	return svcDriver, nil
}

func BuildIngressDriver(ctx context.Context, sbName string, ing create_infra.Ingress, ip create_infra.IngressPaths) (zeus_topology_config_drivers.IngressDriver, error) {
	ingDriver := zeus_topology_config_drivers.IngressDriver{
		Ingress: v1networking.Ingress{
			Spec: v1networking.IngressSpec{
				TLS: []v1networking.IngressTLS{{
					Hosts:      []string{ing.Host},
					SecretName: "tls-secret", // TODO rename
				}},
				Rules: []v1networking.IngressRule{{
					Host: ing.Host,
					IngressRuleValue: v1networking.IngressRuleValue{
						HTTP: &v1networking.HTTPIngressRuleValue{
							Paths: []v1networking.HTTPIngressPath{{
								Path:     "",  // TODO
								PathType: nil, // TODO
								Backend: v1networking.IngressBackend{
									Service:  nil, // todo
									Resource: nil,
								},
							}},
						},
					},
				}},
			},
		},
		Host:         ing.Host,
		NginxAuthURL: ing.AuthServerURL,
	}
	return ingDriver, nil
}

func BuildConfigMapDriver(ctx context.Context, sbName string, configMap create_infra.ConfigMap) (zeus_topology_config_drivers.ConfigMapDriver, error) {
	cmDriver := zeus_topology_config_drivers.ConfigMapDriver{
		ConfigMap: v1.ConfigMap{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Immutable:  nil,
			Data:       nil,
			BinaryData: nil,
		},
		SwapKeys: nil,
	}
	for key, value := range configMap {
		cmDriver.ConfigMap.Data[key] = value
	}
	return cmDriver, nil
}

func LabelBuilder(ctx context.Context) {
	// TODO
}

func BuildContainerDriver(ctx context.Context, sbName string, container create_infra.Container) (v1.Container, error) {
	c := v1.Container{
		Name:    sbName,
		Image:   container.DockerImage.ImageName,
		Command: strings.Split(container.DockerImage.Cmd, ","),
		Args:    strings.Split(container.DockerImage.Args, ","),
		Ports:   []v1.ContainerPort{},
		EnvFrom: nil,
		Env:     nil,
		Resources: v1.ResourceRequirements{
			Limits:   nil,
			Requests: nil,
			Claims:   nil,
		},
		VolumeMounts:    []v1.VolumeMount{},
		ImagePullPolicy: "IfNotPresent",
	}

	if container.IsInitContainer {
		// set init container
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
