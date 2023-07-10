package artemis_eth_rxs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type RxTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *RxTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *RxTestSuite) TestSelect() {
}

func TestRxTestSuite(t *testing.T) {
	suite.Run(t, new(RxTestSuite))
}
