package read_infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyReadActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyReadActionRequestTestSuite) TestReadChart() {
	test.Kns.Namespace = "demo"

	oID, uID := 1667539995595436059, 1667539995591030847
	topID := 1667539995610096
	orgUser := org_users.NewOrgUserWithID(oID, uID)
	tar := TopologyActionReadRequest{
		TopologyActionRequest: base.CreateTopologyActionRequestWithOrgUser("read", orgUser),
		TopologyID:            topID,
	}
	t.Endpoint = "/infra"
	t.AddEndpointHandler(tar.ReadTopology)
	tr := t.PostTopologyRequest(tar, 200)
	t.Require().NotEmpty(tr.Logs)
}

func TestTopologyReadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyReadActionRequestTestSuite))
}
