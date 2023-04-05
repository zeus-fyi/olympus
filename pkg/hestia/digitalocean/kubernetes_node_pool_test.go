package hestia_digitalocean

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DoKubernetesTestSuite struct {
	test_suites_base.TestSuite
}

func (s *DoKubernetesTestSuite) TestGetSizes() {
	s.InitLocalConfigs()

	do := InitDoClient(ctx, "token")
	s.Require().NotNil(do.Client)

	// TODO
}

func (s *DoKubernetesTestSuite) TestCreateNodePool() {
	s.InitLocalConfigs()

	do := InitDoClient(ctx, "token")
	s.Require().NotNil(do.Client)

	// TODO
}

func TestDoKubernetesTestSuite(t *testing.T) {
	suite.Run(t, new(DoKubernetesTestSuite))
}
