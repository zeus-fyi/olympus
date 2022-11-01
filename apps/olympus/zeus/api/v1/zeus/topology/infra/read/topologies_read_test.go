package read_infra

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyReadActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyReadActionRequestTestSuite) TestReadChart() {
	test.Kns.Namespace = "demo"
	bi := base_infra.TopologyInfraActionRequest{
		TopologyActionRequest: base.TopologyActionRequest{
			Action:  "read",
			OrgUser: test.TestOrgUser,
		}}

	tar := TopologyActionReadRequest{
		bi,
		test.TestTopologyID,
	}
	var c echo.Context
	err := tar.ReadTopology(c)
	t.Require().Nil(err)
}

func TestTopologyReadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyReadActionRequestTestSuite))
}
