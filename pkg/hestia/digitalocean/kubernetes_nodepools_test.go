package hestia_digitalocean

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DoKubernetesTestSuite struct {
	test_suites_base.TestSuite
	do DigitalOcean
}

func (s *DoKubernetesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.do = InitDoClient(ctx, s.Tc.DigitalOceanAPIKey)
	s.Require().NotNil(s.do.Client)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}
func (s *DoKubernetesTestSuite) TestGetNodePools() {
	s.InitLocalConfigs()

	nodePools, _, err := s.do.Client.Kubernetes.List(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotEmpty(nodePools)
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
