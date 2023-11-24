package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type K8TestSuite struct {
	test_suites_base.TestSuite
	K K8Util
}

func (s *K8TestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.K = K8Util{}
	s.K.PrintOn = true
	//s.K.ConnectToK8s()
	auth_startup.Ksp.DirIn = "../configs"
	auth_startup.Sp.DirIn = "../configs"
	authCfg := auth_startup.NewDefaultAuthClient(context.Background(), s.Tc.ProdLocalAuthKeysCfg)
	inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(context.Background(), authCfg)
	s.K.ConnectToK8sFromInMemFsCfgPath(inMemFs)
}

func TestK8sTestSuiteTest(t *testing.T) {
	suite.Run(t, new(K8TestSuite))
}
