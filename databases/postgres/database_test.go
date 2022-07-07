package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"

	"bitbucket.org/zeus/eth-indexer/pkg/test_utils"
)

var localConnPG *pgxpool.Pool

type PGTestSuite struct {
	suite.Suite
}

func (s *PGTestSuite) SetupTest() {
	test_utils.InitLocalConfigs()
	localConnPG = test_utils.InitLocalTestDBConn()
}

func (s *PGTestSuite) TestConnPG() {
	ctx := context.Background()
	localConnPG.QueryRow(ctx, "SELECT id FROM validator WHERE id > 0")
}

func TestPGTestSuite(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
