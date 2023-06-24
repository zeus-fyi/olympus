package hestia_ovhcloud

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type OvhCloudTestSuite struct {
	test_suites_base.TestSuite
	o OvhCloud
}

func (s *OvhCloudTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	creds := OvhCloudCreds{
		Region:      OvhUS,
		AppKey:      s.Tc.OvhAppKey,
		AppSecret:   s.Tc.OvhSecretKey,
		ConsumerKey: s.Tc.OvhConsumerKey,
	}
	s.o = InitOvhClient(ctx, creds)
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}
func (s *OvhCloudTestSuite) TestListSizes() {

}

func TestOvhCloudTestSuite(t *testing.T) {
	suite.Run(t, new(OvhCloudTestSuite))
}
