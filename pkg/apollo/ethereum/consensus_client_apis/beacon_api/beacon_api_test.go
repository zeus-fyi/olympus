package beacon_api

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

const disableHighDataAPITests = true

var ctx context.Context

type BeaconAPITestSuite struct {
	test_suites_base.TestSuite
}

func (s *BeaconAPITestSuite) TestGetValidatorsByState() {
	s.SkipTest(disableHighDataAPITests)
	state := "finalized"
	status := "active_ongoing"
	vs, err := GetValidatorsByState(ctx, s.Tc.LocalBeaconConn, state, status)

	s.Require().Nil(err)
	file, _ := json.Marshal(vs)

	_ = os.WriteFile("validators.json", file, 0644)
}

func (s *BeaconAPITestSuite) TestGetValidatorsByStateFilter() {
	s.T().Parallel()
	state := "head"
	valIndexes := []string{"242521", "67596"}
	encodedURLparams := string_utils.UrlExplicitEncodeQueryParamList("id", valIndexes...)
	vs, err := GetValidatorsBalancesByStateFilter(ctx, s.Tc.LocalBeaconConn, state, encodedURLparams)
	s.Require().Nil(err)
	s.Assert().Len(vs.Data, 2)
	s.Assert().Equal(false, vs.ExecutionOptimistic)
}

func (s *BeaconAPITestSuite) TestGetBlockByID() {
	s.T().Parallel()
	r := GetBlockByID(ctx, s.Tc.LocalBeaconConn, "head")
	s.Require().Nil(r.Err)
}

func TestBeaconAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconAPITestSuite))
}
