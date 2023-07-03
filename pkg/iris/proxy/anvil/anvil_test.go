package proxy_anvil

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type AnvilTestSuite struct {
	test_suites_base.TestSuite
}

func (t *AnvilTestSuite) SetupTest() {
	t.InitLocalConfigs()
	InitAnvilProxy()
	t.Assert().NotNil(SessionLocker.LFU)
}

func (t *AnvilTestSuite) TestSessionLocker() {
	r, err := SessionLocker.GetNextAvailableRouteAndAssignToSession("test")
	t.Assert().Nil(err)
	t.Assert().NotNil(r)
}

func TestAnvilTestSuite(t *testing.T) {
	suite.Run(t, new(AnvilTestSuite))
}
