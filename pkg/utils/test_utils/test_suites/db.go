package test_suites

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"

	"github.com/zeus-fyi/olympus/configs"
)

type PGTestSuite struct {
	base.TestSuite
	Pg      apps.Db
	LocalDB bool
}

func (s *PGTestSuite) SetupTest() {
	s.SetupPGConn()
}

func (s *PGTestSuite) SetupLocalPGConn() {
	s.Tc = configs.InitStagingConfigs()
	s.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
	s.LocalDB = true
}

func (s *PGTestSuite) SetupStagingPGConn() {
	s.Tc = configs.InitStagingConfigs()
	log.Warn().Msg("WARNING: staging db connection")
	s.Pg.InitPG(context.Background(), s.Tc.StagingDbPgconn)
	s.LocalDB = false
}

func (s *PGTestSuite) SetupProductionPGConn() {
	s.Tc = configs.InitProductionConfigs()
	log.Warn().Msg("WARNING: production db connection")
	s.Pg.InitPG(context.Background(), s.Tc.ProdDbPgconn)
	s.LocalDB = false
}

func (s *PGTestSuite) SetupPGConn() {
	s.Tc = configs.InitLocalTestConfigs()
	switch s.Tc.Env {
	case "local":
		s.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
		s.LocalDB = true
	case "staging":
		// staging
		s.Pg.InitPG(context.Background(), s.Tc.StagingDbPgconn)
		s.LocalDB = false
	case "production":
		log.Warn().Msg("WARNING: production db connection")
		s.Pg.InitPG(context.Background(), s.Tc.ProdDbPgconn)
		s.LocalDB = false
		return
	default:
		s.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
		s.LocalDB = true
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
		_, err := apps.Pg.Exec(ctx, query)
		log.Err(err).Interface("cleanupDb: %s", tableName)
		if err != nil {
			panic(err)
		}
	}
}

func TestPGTestSuite(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
