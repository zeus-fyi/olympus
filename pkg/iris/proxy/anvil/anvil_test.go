package proxy_anvil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type AnvilTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (t *AnvilTestSuite) SetupTest() {
	t.InitLocalConfigs()
	InitAnvilProxy()
	SessionLocker.LockDefaultTime = time.Second * 1
	t.Require().NotNil(SessionLocker.PriorityQueue)
}

func (t *AnvilTestSuite) TestSessionLocker() {
	r, err := SessionLocker.GetSessionLockedRoute(ctx, "test1")
	t.Require().Nil(err)
	t.Assert().NotNil(r)
	fmt.Println(r)
	r2, err := SessionLocker.GetSessionLockedRoute(ctx, "test2")
	t.Require().Nil(err)
	t.Assert().NotNil(r2)
	fmt.Println(r2)
	r2again, err := SessionLocker.GetSessionLockedRoute(ctx, "test2")
	t.Assert().Nil(err)
	t.Require().Equal(r2again, r2)
	time.Sleep(time.Second * 1)

	now := time.Now()
	for i := 0; i < 20; i++ {
		nr, err := SessionLocker.GetSessionLockedRoute(ctx, fmt.Sprintf("%d", i))
		t.Assert().Nil(err)
		t.Assert().NotNil(nr)
		fmt.Println(nr)
	}
	fmt.Println("no removals", time.Since(now))
	now = time.Now()

	for i := 0; i < 20; i++ {
		nr, err := SessionLocker.GetSessionLockedRoute(ctx, fmt.Sprintf("%d", i))
		t.Assert().Nil(err)
		t.Assert().NotNil(nr)
		SessionLocker.RemoveSessionLockedRoute(fmt.Sprintf("%d", i))
	}
	fmt.Println("with removals", time.Since(now))

}
func TestAnvilTestSuite(t *testing.T) {
	suite.Run(t, new(AnvilTestSuite))
}
