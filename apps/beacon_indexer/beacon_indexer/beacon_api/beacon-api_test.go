package beacon_api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_api/api_types"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

const disableHighDataAPITests = false

var ctx context.Context

type BeaconAPITestSuite struct {
	base.BaseTestSuite
}

func (s *BeaconAPITestSuite) TestGetValidatorsByState() {
	s.SkipTest(disableHighDataAPITests)
	state := "finalized"

	r := GetValidatorsByState(ctx, s.Tc.BeaconNodeInfura, state)
	s.Require().Nil(r.Err)
	var vs api_types.ValidatorsStateBeacon
	err := json.Unmarshal(r.BodyBytes, &vs)
	s.Require().Nil(err)
	file, _ := json.Marshal(vs)

	_ = ioutil.WriteFile("validators.json", file, 0644)
}

func (s *BeaconAPITestSuite) TestGetValidatorsByStateFilter() {
	s.T().Parallel()
	state := "head"
	valIndexes := []string{"242521", "67596"}
	encodedURLparams := string_utils.UrlEncodeQueryParamList("", valIndexes...)
	r := GetValidatorsBalancesByStateFilter(ctx, s.Tc.BeaconNodeInfura, state, encodedURLparams)
	s.Require().Nil(r.Err)

	var vb api_types.ValidatorBalances
	err := json.Unmarshal(r.BodyBytes, &vb)
	s.Require().Nil(err)
	s.Assert().Len(vb.Data, 2)
	s.Assert().Equal(false, vb.ExecutionOptimistic)
}

func (s *BeaconAPITestSuite) TestGetBlockByID() {
	s.T().Parallel()
	r := GetBlockByID(ctx, s.Tc.BeaconNodeInfura, "head")
	s.Require().Nil(r.Err)
}

func TestBeaconAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconAPITestSuite))
}
