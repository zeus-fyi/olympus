package nodes

import (
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type NodesActionsRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *NodesActionsRequestTestSuite) TestEndToEnd() {
	//t.InitLocalConfigs()
	//
	//t.Eg.POST("/infra/create", CreateTopologyInfraActionRequestHandler)
	//
	//
	//start := make(chan struct{}, 1)
	//go func() {
	//	close(start)
	//	_ = t.E.Start(":9010")
	//}()
	//
	//ctx := context.Background()

}

func TestNodesActionsRequestTestSuite(t *testing.T) {
	suite.Run(t, new(NodesActionsRequestTestSuite))
}
