package artemis_trading_auxiliary

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
)

type ArtemisAuxillaryTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
	acc        accounts.Account
	goerliNode string
}

var ctx = context.Background()

func (t *ArtemisAuxillaryTestSuite) SetupTest() {
	t.S3SecretsManagerTestSuite.SetupTest()
	t.goerliNode = t.Tc.GoerliNodeUrl
	tc := t.S3SecretsManagerTestSuite.Tc
	athena.AthenaS3Manager = t.S3SecretsManagerTestSuite.S3
	apps.Pg.InitPG(ctx, tc.ProdLocalDbPgconn)
	age := encryption.NewAge(tc.LocalAgePkey, tc.LocalAgePubkey)
	t.acc = artemis_trading_cache.InitAccount(ctx, age)
}

func TestArtemisAuxiliaryTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisAuxillaryTestSuite))
}
