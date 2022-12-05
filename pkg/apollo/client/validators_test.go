package apollo_client

import (
	"github.com/zeus-fyi/olympus/pkg/apollo/client/req_types"
)

func (t *ApolloClientTestSuite) TestValidatorStatuses() {
	rr := req_types.ValidatorsRequest{ValidatorIndexes: []int{1, 483925}}
	resp, err := t.ApolloTestClient.ValidatorStatuses(ctx, rr)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *ApolloClientTestSuite) TestValidatorBalances() {
	rr := req_types.ValidatorBalancesRequest{ValidatorIndexes: []int{1, 2}, LowerEpoch: 164000, HigherEpoch: 164010}
	resp, err := t.ApolloTestClient.ValidatorBalances(ctx, rr)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *ApolloClientTestSuite) TestValidatorBalanceSums() {
	rr := req_types.ValidatorBalancesRequest{ValidatorIndexes: []int{1, 2}, LowerEpoch: 164000, HigherEpoch: 164010}
	resp, err := t.ApolloTestClient.ValidatorBalanceSums(ctx, rr)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}
