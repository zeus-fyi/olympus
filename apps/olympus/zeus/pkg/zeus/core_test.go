package zeus

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type ZeusAppCoreTestSuite struct {
	test_suites_base.TestSuite
}

func (s *ZeusAppCoreTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	authKeysCfg := s.Tc.ProdLocalAuthKeysCfg
	KeysCfg = auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
}

var ctx = context.Background()

func (s *ZeusAppCoreTestSuite) TestVerifyClusterAuthAndGetKubeCfg() {
	cctx := zeus_common_types.CloudCtxNs{
		ClusterCfgStrID: "1707274735429812000",
		CloudProvider:   "aws",
		Region:          "us-east-2",
		Context:         "zeus-eks-us-east-2",
		Namespace:       "anvil-2f7c37f6",
		Alias:           "anvil-2f7c37f6",
	}
	k, err := VerifyClusterAuthAndGetKubeCfg(ctx, s.Ou, cctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(k)

	ns, err := k.GetNamespaces(ctx, cctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(ns)

	for _, n := range ns.Items {
		if strings.HasPrefix(n.Name, "anvil") {
			cctx.Namespace = n.Name
			err = k.DeleteNamespace(ctx, cctx)
			s.Require().Nil(err)
		}
	}
}

func (s *ZeusAppCoreTestSuite) TestClusterAuthKube() {
	p := authorized_clusters.K8sClusterConfig{
		ExtConfigStrID: "",
		ExtConfigID:    0,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			ClusterCfgStrID: "1707274735429812000",
			CloudProvider:   "aws",
			Region:          "us-east-2",
			Context:         "zeus-eks-us-east-2",
			Namespace:       "anvil-2f7c37f6",
			Alias:           "anvil-2f7c37f6",
		},
		IsActive: false,
		IsPublic: false,
	}

	resp, err := authorized_clusters.SelectAuthedClusterByRouteAndOrgID(ctx, s.Ou, p.CloudCtxNs)
	s.Require().Nil(err)
	s.Require().NotNil(resp)
	s.Require().NotEmpty(resp)
	_, err = GetKubeConfig(ctx, s.Ou, p)
	s.Require().Nil(err)
}

func TestZeusAppCoreTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusAppCoreTestSuite))
}
