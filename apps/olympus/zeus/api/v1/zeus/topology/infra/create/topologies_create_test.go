package create_infra

import (
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestCreateTopology() {
	name := fmt.Sprintf("random_%d", t.Ts.UnixTimeStampNow())

	topologyActionRequest := TopologyActionCreateRequest{
		TopologyInfraActionRequest: base_infra.TopologyInfraActionRequest{},
		TopologyCreateRequest:      TopologyCreateRequest{Name: name},
	}

	var c echo.Context
	err := topologyActionRequest.CreateTopology(c)
	t.Require().Nil(err)
}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
