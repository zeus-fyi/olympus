package test_suites

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/databases/postgres"

	"github.com/zeus-fyi/olympus/configs"
)

type PGTestSuite struct {
	BaseTestSuite
	Pg postgres.Db
}

func (s *PGTestSuite) SetupTest() {
	s.Tc = configs.InitLocalTestConfigs()
	s.Pg.InitPG(context.Background(), s.Tc.TEST_DB_PGCONN)
}

func TestPGTestSuite(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
