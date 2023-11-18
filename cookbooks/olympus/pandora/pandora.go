package pandora

import (
	"context"
	"fmt"

	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	zeus_cluster_config_drivers "github.com/zeus-fyi/zeus/zeus/cluster_config_drivers"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/zk8s_templates"
	zeus_client "github.com/zeus-fyi/zeus/zeus/z_client"
)

var (
	pandoraClusterDefinition = zeus_cluster_config_drivers.ClusterDefinition{
		ClusterClassName: "pandora",
		ComponentBases:   pandoraComponentBases,
	}
	pandoraComponentBases = map[string]zeus_cluster_config_drivers.ComponentBaseDefinition{
		"pandora": pandoraComponentBase,
	}
	pandoraComponentBase = zeus_cluster_config_drivers.ComponentBaseDefinition{
		SkeletonBases: map[string]zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
			"pandora": pandoraSkeletonBaseConfig,
		},
	}
	pandoraSkeletonBaseConfig = zeus_cluster_config_drivers.ClusterSkeletonBaseDefinition{
		SkeletonBaseNameChartPath: pandoraChartPath,
	}
)

var pandoraChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/pandora/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "pandora", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
}

func CreatePandora(ctx context.Context, zc zeus_client.ZeusClient, createClass bool) error {

	docusaurusTemplate := "pandora"
	wd := zeus_cluster_config_drivers.WorkloadDefinition{
		WorkloadName: docusaurusTemplate,
		ReplicaCount: 1,
		Containers: zk8s_templates.Containers{
			docusaurusTemplate: zk8s_templates.Container{
				IsInitContainer: false,
				ImagePullPolicy: "Always",
				DockerImage: zk8s_templates.DockerImage{
					ImageName: "registry.digitalocean.com/zeus-fyi/pandora:latest",
					ResourceRequirements: zk8s_templates.ResourceRequirements{
						CPU:    "100m",
						Memory: "500Mi",
					},
					Ports: []zk8s_templates.Port{
						{
							Name:               "http",
							Number:             "8000",
							Protocol:           "TCP",
							IngressEnabledPort: false,
							ProbeSettings: zk8s_templates.ProbeSettings{
								UseForLivenessProbe:  true,
								UseForReadinessProbe: true,
								UseTcpSocket:         true,
							},
						},
					},
				},
			},
		},
	}
	cd, err := zeus_cluster_config_drivers.GenerateDeploymentCluster(ctx, wd)
	if err != nil {
		return err
	}

	prt, err := zeus_cluster_config_drivers.PreviewTemplateGeneration(ctx, cd)
	if err != nil {
		return err
	}

	if createClass {
		gcd := zeus_cluster_config_drivers.CreateGeneratedClusterClassCreationRequest(cd)
		fmt.Println(gcd)
		err = gcd.CreateClusterClassDefinitions(ctx, zc)
		if err != nil {
			return err
		}
	}

	_, err = prt.UploadChartsFromClusterDefinition(ctx, zc, true)
	if err != nil {
		return err
	}
	return nil
}
