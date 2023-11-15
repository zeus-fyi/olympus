package nodes

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var ctx = context.Background()

type NodesActionsRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *NodesActionsRequestTestSuite) TestNodes() {
	t.InitLocalConfigs()
	t.Eg.POST("/nodes", NodeActionsRequestHandler)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	req := ActionRequest{
		Action: "list",
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "ovh",
			Region:        "us-west-or-1",
			Context:       "kubernetes-admin@zeusfyi",
		},
	}

	nl := zeus_core.ClusterNodesAudit{}
	resp, err := t.ZeusClient.R().
		SetResult(&nl).
		SetBody(&req).
		Post("/v1/nodes")
	t.Require().Nil(err)
	t.Require().Equal(http.StatusOK, resp.StatusCode())
	t.Assert().NotEmpty(resp)
}

func TestNodesActionsRequestTestSuite(t *testing.T) {
	suite.Run(t, new(NodesActionsRequestTestSuite))
}
