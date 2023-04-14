package zeus_core

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type NodesTestSuite struct {
	K8TestSuite
}

func (t *NodesTestSuite) TestGetNodes() {
	ctx := context.Background()
	//var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: ""}
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "nyc1", Context: "do-nyc1-do-nyc1-zeus-demo", Namespace: ""}

	nodes, err := t.K.GetNodes(ctx, kns)
	t.Require().Nil(err)
	t.Require().NotEmpty(nodes)
}

func (t *NodesTestSuite) TestGetNodesByLabel() {
	ctx := context.Background()
	//var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: ""}
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "nyc1", Context: "do-nyc1-do-nyc1-zeus-demo", Namespace: ""}
	// org -> 1679515557647002001
	m := make(map[string]string)
	m["org"] = "1679515557647002001"
	m["app"] = "ethereumEphemeralBeacons"

	var labelString string
	for key, value := range m {
		labelString += key + "=" + value + ","
	}
	labelString = strings.TrimSuffix(labelString, ",")
	nl, err := t.K.GetNodesAuditByLabel(ctx, kns, labelString)
	t.Require().Nil(err)
	t.Require().NotEmpty(nl)
}

func TestNodesTestSuite(t *testing.T) {
	suite.Run(t, new(NodesTestSuite))
}
