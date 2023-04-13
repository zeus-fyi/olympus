package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type NodesTestSuite struct {
	K8TestSuite
}

func (t *NodesTestSuite) TestGetNodes() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: ""}

	nodes, err := t.K.GetNodes(ctx, kns)
	t.Require().Nil(err)
	t.Require().NotEmpty(nodes)
}

func TestNodesTestSuite(t *testing.T) {
	suite.Run(t, new(NodesTestSuite))
}
