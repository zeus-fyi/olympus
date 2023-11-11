package artemis_rawdawg_contract

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

var ctx = context.Background()

type ArtemisTradingContractsTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	LocalEnv web3_client.Web3Client
}

func (s *ArtemisTradingContractsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}

func TestArtemisTradingContractsTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradingContractsTestSuite))
}
