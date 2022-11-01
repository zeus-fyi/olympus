package create_infra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestCreateTopology() {
	name := fmt.Sprintf("random_%d", t.Ts.UnixTimeStampNow())

	tar := TopologyActionCreateRequest{
		TopologyInfraActionRequest: base_infra.TopologyInfraActionRequest{},
		TopologyCreateRequest:      TopologyCreateRequest{Name: name},
	}

	t.Endpoint = "/infra"
	t.AddEndpointHandler(tar.CreateTopology)
	tr := t.PostTopologyRequest(tar, 200)
	t.Require().NotEmpty(tr.Logs)
}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
