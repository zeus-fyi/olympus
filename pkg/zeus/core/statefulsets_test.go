package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type StatefulSetsTestSuite struct {
	K8TestSuite
}

func (s *StatefulSetsTestSuite) TestRolloutRestartDeployment() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "zeus"}
	err := s.K.RolloutRestartStatefulSet(ctx, kns, "tbd", nil)
	s.Require().Nil(err)
}

func TestStatefulSetsTestSuite(t *testing.T) {
	suite.Run(t, new(StatefulSetsTestSuite))
}
