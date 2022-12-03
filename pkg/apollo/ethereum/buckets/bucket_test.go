package apollo_buckets

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/consensus_client_apis/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_s3"
)

type ApolloBucketsTestSuite struct {
	test_suites_s3.S3TestSuite
}

func (s *ApolloBucketsTestSuite) TestGetBuckets() {
	ctx := context.Background()
	tc := configs.InitLocalTestConfigs()
	ApolloS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, tc.DevAuthKeysCfg)

	ba, err := DownloadBalancesAtEpoch(ctx, "validator-balance-epoch-164046")
	s.Require().Nil(err)
	vbe := beacon_api.ValidatorBalances{}
	err = json.Unmarshal(ba, &vbe)
	s.Require().Nil(err)
	s.Require().NotEmpty(vbe)
}

func TestApolloBucketsTestSuite(t *testing.T) {
	suite.Run(t, new(ApolloBucketsTestSuite))
}
