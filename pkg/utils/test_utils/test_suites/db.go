package test_suites

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"

	"github.com/zeus-fyi/olympus/configs"
)

type PGTestSuite struct {
	base.TestSuite
	Pg postgres.Db
}

func (s *PGTestSuite) SetupTest() {
	s.Tc = configs.InitLocalTestConfigs()
	if len(s.Tc.LocalDbPgconn) > 0 {
		// local
		s.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
	} else {
		// staging
		s.Pg.InitPG(context.Background(), s.Tc.StagingDbPgconn)
	}
}

func TestPGTestSuite(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
