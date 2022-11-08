package topology_worker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type TopologyWorkerTestSuite struct {
	test_suites.TemporalTestSuite
}

func (s *TopologyWorkerTestSuite) SetupTest() {
}

func (s *TopologyWorkerTestSuite) TestCreateWorker() {
	tc := configs.InitLocalTestConfigs()
	fmt.Println("prod local")

	w, err := InitTopologyWorker(tc.ProdLocalTemporalAuth)
	s.Assert().Nil(err)
	s.Assert().NotEmpty(w)
}

func TestTopologyWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyWorkerTestSuite))
}
