package artemis_etherscan

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

var ctx = context.Background()

type EtherscanTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	es Etherscan
}

func (s *EtherscanTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	s.es = NewMainnetEtherscanClient(s.Tc.EtherScanAPIKey)
}

func (s *EtherscanTestSuite) TestFetchErc20Info() {

}

func TestEtherscanTestSuite(t *testing.T) {
	suite.Run(t, new(EtherscanTestSuite))
}
