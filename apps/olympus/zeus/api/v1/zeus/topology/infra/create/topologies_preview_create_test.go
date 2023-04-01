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

	// TODO
	req := Cluster{
		ClusterName:     "clusterTest",
		ComponentBases:  make(map[string]SkeletonBases),
		IngressSettings: Ingress{},
		IngressPaths:    IngressPaths{},
	}

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
