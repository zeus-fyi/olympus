package beacon_api

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
	r := GetValidatorsByState(ctx, s.Tc.LocalBeaconConn, state, status)
	s.Require().Nil(r.Err)
	var vs ValidatorsStateBeacon
	err := json.Unmarshal(r.BodyBytes, &vs)
	s.Require().Nil(err)
	file, _ := json.Marshal(vs)

	_ = ioutil.WriteFile("validators.json", file, 0644)
}

func (s *BeaconAPITestSuite) TestGetValidatorsByStateFilter() {
	s.T().Parallel()
	state := "head"
	valIndexes := []string{"242521", "67596"}
	encodedURLparams := string_utils.UrlExplicitEncodeQueryParamList("id", valIndexes...)
	r := GetValidatorsBalancesByStateFilter(ctx, s.Tc.LocalBeaconConn, state, encodedURLparams)
	s.Require().Nil(r.Err)

	var vb ValidatorBalances
	err := json.Unmarshal(r.BodyBytes, &vb)
	s.Require().Nil(err)
	s.Assert().Len(vb.Data, 2)
	s.Assert().Equal(false, vb.ExecutionOptimistic)
}

func (s *BeaconAPITestSuite) TestGetBlockByID() {
	s.T().Parallel()
	r := GetBlockByID(ctx, s.Tc.LocalBeaconConn, "head")
	s.Require().Nil(r.Err)
}

func TestBeaconAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconAPITestSuite))
}
