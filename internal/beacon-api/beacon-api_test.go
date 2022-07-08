package beacon_api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/internal/beacon-api/api_types"
	"github.com/zeus-fyi/olympus/pkg/test_utils/test_suites"
)

const disableHighDataAPITests = true

type BeaconAPITestSuite struct {
	test_suites.BaseTestSuite
}

func (s *BeaconAPITestSuite) TestGetValidatorsByState() {
	s.SkipTest(disableHighDataAPITests)
	state := "finalized"

	r := GetValidatorsByState(s.Tc.BEACON_NODE_INFURA, state)
	s.Require().Nil(r.Err)
}

func (s *BeaconAPITestSuite) TestGetValidatorsByStateFilter() {
	s.T().Parallel()
	state := "head"
	valIndexes := []string{"242521", "67596"}
	r := GetValidatorsBalancesByStateFilter(s.Tc.BEACON_NODE_INFURA, state, valIndexes...)
	s.Require().Nil(r.Err)

	var vb api_types.ValidatorBalances
	err := json.Unmarshal(r.BodyBytes, &vb)
	s.Require().Nil(err)
	s.Assert().Len(vb.Data, 2)
	s.Assert().Equal(false, vb.ExecutionOptimistic)
}

func (s *BeaconAPITestSuite) TestGetBlockByID() {
	s.T().Parallel()
	r := GetBlockByID(s.Tc.BEACON_NODE_INFURA, "head")
	s.Require().Nil(r.Err)
}

func TestBeaconAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconAPITestSuite))
}
