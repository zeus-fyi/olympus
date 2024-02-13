package deploy_topology_activities_create_setup

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type HelmDeployTestSuite struct {
	test_suites_base.TestSuite
}

func (s *HelmDeployTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *HelmDeployTestSuite) TestDeployHelmChart() {

	DeployHelmChart()
}

func TestHelmDeployTestSuite(t *testing.T) {
	suite.Run(t, new(HelmDeployTestSuite))
}
