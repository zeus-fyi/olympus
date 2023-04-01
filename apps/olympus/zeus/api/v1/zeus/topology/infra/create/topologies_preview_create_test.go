package create_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyPreviewCreateClassRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

const previewEndpoint = "/infra/preview/create"

func (t *TopologyPreviewCreateClassRequestTestSuite) TestGeneratePreview() {
	t.InitLocalConfigs()

	t.Eg.POST("/infra/preview/create", PreviewCreateTopologyInfraActionRequestHandler)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

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

	var jsonResp any
	resp, err := t.ZeusClient.R().
		SetResult(&jsonResp).
		SetBody(&req).
		Post(previewEndpoint)

	t.Require().Nil(err)
	t.ZeusClient.PrintRespJson(resp.Body())
}

func TestTopologyPreviewCreateClassRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyPreviewCreateClassRequestTestSuite))
}
