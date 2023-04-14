package nodes

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
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
		TopologyKubeCtxNs: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "",
				Region:        "",
				Context:       "do-nyc1-do-nyc1-zeus-demo",
				Namespace:     "",
				Env:           "",
			},
		},
		Action: "list",
	}

	nl := []zeus_core.NodeAudit{}
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
