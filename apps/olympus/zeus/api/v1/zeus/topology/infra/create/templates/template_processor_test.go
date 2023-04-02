package zeus_templates

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TemplateProcessorTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *TemplateProcessorTestSuite) TestGeneratePreview() {
	t.InitLocalConfigs()

	ctx := context.Background()

	req := Cluster{
		ClusterName:     "avaxNodeTest",
		ComponentBases:  make(map[string]SkeletonBases),
		IngressSettings: Ingress{},
		IngressPaths:    IngressPaths{},
	}

	sb := SkeletonBase{
		AddStatefulSet:    true,
		AddDeployment:     false,
		AddConfigMap:      false,
		AddService:        true,
		AddIngress:        false,
		AddServiceMonitor: false,
		ConfigMap:         ConfigMap{},
		StatefulSet: StatefulSet{
			ReplicaCount: 1,
			PVCTemplates: []PVCTemplate{{
				Name:               "avax-client-storage",
				AccessMode:         "ReadWriteOnce",
				StorageSizeRequest: "2Ti",
			}},
		},
		Containers: make(map[string]Container),
	}
	c := Container{
		IsInitContainer: false,
		DockerImage: DockerImage{
			ImageName: "avaplatform/avalanchego:v1.9.10",
			Cmd:       "/bin/sh",
			Args:      "-c,/scripts/start.sh",
			ResourceRequirements: ResourceRequirements{
				CPU:    "6",
				Memory: "12Gi",
			},
			Ports: []Port{
				{
					Name:               "p2p-tcp",
					Number:             "9651",
					Protocol:           "TCP",
					IngressEnabledPort: false,
				}, {
					Name:               "http-api",
					Number:             "9650",
					Protocol:           "TCP",
					IngressEnabledPort: true,
				}, {
					Name:               "metrics",
					Number:             "9090",
					Protocol:           "TCP",
					IngressEnabledPort: false,
				},
			},
			VolumeMounts: []VolumeMount{{
				Name:      "avax-client-storage",
				MountPath: "/data",
			}},
		},
	}

	sb.Containers["avax-client"] = c
	req.ComponentBases["avaxClients"] = make(map[string]SkeletonBase)
	req.ComponentBases["avaxClients"]["avaxClients"] = sb

	cd := PreviewTemplateGeneration(ctx, req)
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := cd.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := cd.GenerateSkeletonBaseCharts()
	t.Assert().NoError(err)
	t.Assert().NotEmpty(sbDefs)
}

func TestTemplateProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateProcessorTestSuite))
}
