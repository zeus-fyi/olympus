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
		fmt.Println(nr)
		t.Assert().Nil(err)
		t.Assert().NotNil(nr)
		if i == 0 {
			fmt.Println("locked", nr)
		}
		route, err := SessionLocker.GetSessionLockedRoute(ctx, fmt.Sprintf("%d", 0))
		t.Assert().Nil(err)
		t.Assert().NotNil(route)
	}
	fmt.Println("locking one", time.Since(now))

	now = time.Now()
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
		SessionLocker.RemoveSessionLockedRoute(ctx, fmt.Sprintf("%d", i))
	}
	fmt.Println("with removals", time.Since(now))
	/*
		no removals 4.90930125s
		with removals 1.999989417s
	*/
}
func TestAnvilTestSuite(t *testing.T) {
	suite.Run(t, new(AnvilTestSuite))
}
