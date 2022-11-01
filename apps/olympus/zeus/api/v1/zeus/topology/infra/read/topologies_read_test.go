package read_infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyReadActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyReadActionRequestTestSuite) TestReadChart() {
	test.Kns.Namespace = "demo"
	tar := TopologyActionReadRequest{
		TopologyActionRequest: base.CreateTopologyActionRequestWithOrgUser("read", test.TestOrgUser),
		TopologyID:            test.TestTopologyID,
	}
	t.Endpoint = "/infra"
	t.AddEndpointHandler(tar.ReadTopology)
	tr := t.PostTopologyRequest(tar, 200)
	t.Require().NotEmpty(tr.Logs)
}

func TestTopologyReadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyReadActionRequestTestSuite))
}
