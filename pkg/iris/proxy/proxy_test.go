package iris_proxy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type IrisProxyTestSuite struct {
	test_suites_base.TestSuite
}

func (s *IrisProxyTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *IrisProxyTestSuite) TestProxy() {

}

func TestIrisProxyTestSuite(t *testing.T) {
	suite.Run(t, new(IrisProxyTestSuite))
}
