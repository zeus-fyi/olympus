package proxy_anvil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type AnvilTestSuite struct {
	test_suites_base.TestSuite
}

func (t *AnvilTestSuite) SetupTest() {
	t.InitLocalConfigs()
	InitAnvilProxy()
	SessionLocker.LockDefaultTime = time.Second * 1
	t.Assert().NotNil(SessionLocker.LFU)
}

func (t *AnvilTestSuite) TestSessionLocker() {
	r, err := SessionLocker.GetNextAvailableRouteAndAssignToSession("test1")
	t.Assert().Nil(err)
	t.Assert().NotNil(r)
	r2, err := SessionLocker.GetNextAvailableRouteAndAssignToSession("test2")
	t.Assert().Nil(err)
	t.Assert().NotNil(r2)

	_, err = SessionLocker.GetSessionLockedRoute("test2")
	t.Assert().Nil(err)
	lfKey, lfVal := SessionLocker.LFU.GetLeastFrequentValue()
	t.Assert().NotNil(lfKey)
	t.Assert().Equal("test1", lfKey)
	fmt.Println(lfVal)

	time.Sleep(time.Second * 1)
	nr, err := SessionLocker.GetNextAvailableRouteAndAssignToSession("test3")
	t.Assert().Nil(err)
	t.Assert().NotNil(nr)
}
func TestAnvilTestSuite(t *testing.T) {
	suite.Run(t, new(AnvilTestSuite))
}
