package test_suites

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
)

type PGTestSuite struct {
	BaseTestSuite
	P *pgxpool.Pool
}

func (s *PGTestSuite) SetupTest() {
	s.Tc = configs.InitLocalTestConfigs()
	conn, err := pgxpool.Connect(context.Background(), s.Tc.TEST_DB_PGCONN)
	if err != nil {
		panic(err)
	}
	s.P = conn
}

func TestPGTestSuite(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
