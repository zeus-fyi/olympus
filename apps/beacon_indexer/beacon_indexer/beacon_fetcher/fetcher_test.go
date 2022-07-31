package beacon_fetcher

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type FetcherBaseTestSuite struct {
	test_suites.PGTestSuite
	P *pgxpool.Pool
}

var batchSize = 10

func (s *FetcherBaseTestSuite) TestBeaconFindNewAndMissingValidatorIndexes() {
	var f BeaconFetcher
	f.NodeEndpoint = s.Tc.LocalBeaconConn
	ctx := context.Background()
	err := f.BeaconFindNewAndMissingValidatorIndexes(ctx, batchSize)
	s.Require().Nil(err)
	s.Assert().Len(f.Validators.Validators, batchSize)
}

func (s *FetcherBaseTestSuite) TestFindAndQueryAndUpdateValidatorBalances() {
	var f BeaconFetcher
	f.NodeEndpoint = s.Tc.LocalBeaconConn
	ctx := context.Background()
	err := f.FindAndQueryAndUpdateValidatorBalances(ctx, batchSize)
	s.Require().Nil(err)
}

func (s *FetcherBaseTestSuite) TestBeaconUpdateValidatorStates() {
	var f BeaconFetcher
	f.NodeEndpoint = s.Tc.LocalBeaconConn
	ctx := context.Background()
	err := f.BeaconUpdateValidatorStates(ctx, batchSize)
	s.Require().Nil(err)
	s.Assert().Len(f.Validators.Validators, batchSize)
}

func TestFetcherBaseTestSuite(t *testing.T) {
	suite.Run(t, new(FetcherBaseTestSuite))
}
