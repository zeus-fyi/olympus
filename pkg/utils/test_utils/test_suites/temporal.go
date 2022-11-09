package test_suites

import (
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_client "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type TemporalTestSuite struct {
	base.TestSuite

	TemporalAuthCfg auth.TemporalAuth
	Temporal        temporal_client.TemporalClient
	Redis           *redis.Client
	PG              apps.Db
	PGTest          PGTestSuite
	RedisTest       RedisTestSuite
}

func (t *TemporalTestSuite) GetTemporalDevAuthCfg() {
	certPath := "./zeus.fyi/ca.pem"
	pemPath := "./zeus.fyi/ca.key"
	namespace := t.Tc.DevTemporalNs
	hostPort := t.Tc.DevTemporalHostPort
	auth := auth.TemporalAuth{
		ClientCertPath:   certPath,
		ClientPEMKeyPath: pemPath,
		Namespace:        namespace,
		HostPort:         hostPort,
	}
	t.TemporalAuthCfg = auth
}

func (t *TemporalTestSuite) SetupTemporalWithPG() {
	t.InitLocalConfigs()
	t.GetTemporalDevAuthCfg()
	client, err := temporal_client.NewTemporalClient(t.TemporalAuthCfg)
	t.Require().Nil(err)
	t.Temporal = client
	t.PGTest.SetupLocalPGConn()
}

func (t *TemporalTestSuite) SetupTest() {
	t.SetupTemporalWithPG()
}

func TestTemporalTestSuite(t *testing.T) {
	suite.Run(t, new(TemporalTestSuite))
}
