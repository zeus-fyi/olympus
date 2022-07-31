package test_suites

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"

	"github.com/zeus-fyi/olympus/configs"
)

type PGTestSuite struct {
	base.TestSuite
	Pg      postgres.Db
	LocalDB bool
}

func (s *PGTestSuite) SetupTest() {
	s.SetupPGConn()
}

func (s *PGTestSuite) SetupPGConn() {
	s.Tc = configs.InitLocalTestConfigs()
	if len(s.Tc.LocalDbPgconn) > 0 {
		// local
		s.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
		s.LocalDB = true
	} else {
		// staging
		s.Pg.InitPG(context.Background(), s.Tc.StagingDbPgconn)
	}
}

func (s *PGTestSuite) CleanupDb(ctx context.Context, tablesToCleanup []string) {
	if s.LocalDB != true {
		log.Info().Msg("not a local database, CleanupDb should only be used on a local database")
		return
	}

	switch s.Tc.Env {
	case "local":
	case "staging":
		log.Info().Msg("not a local database, CleanupDb should only be used on a local database")
		return
	case "production":
		log.Info().Msg("not a local database, CleanupDb should only be used on a local database")
		return
	default:
	}

	for _, tableName := range tablesToCleanup {
		query := fmt.Sprintf(`DELETE FROM %s WHERE %s`, tableName, "true")
		_, err := postgres.Pg.Exec(ctx, query)
		log.Err(err).Interface("cleanupDb: %s", tableName)
		if err != nil {
			panic(err)
		}
	}
}

func TestPGTestSuite(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
