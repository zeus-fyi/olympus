package topology_worker

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type TopologyWorkerTestSuite struct {
	test_suites.TemporalTestSuite
}

func (s *TopologyWorkerTestSuite) TestCreateWorker() {
	w, err := InitTopologyWorker(s.TemporalAuthCfg)
	s.Assert().Nil(err)
	s.Assert().NotEmpty(w)
}

func TestTopologyWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyWorkerTestSuite))
}
