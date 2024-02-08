package zeus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type ZeusAppCoreTestSuite struct {
	test_suites_base.TestSuite
}

func (s *ZeusAppCoreTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *ZeusAppCoreTestSuite) TestUnGzipIntoMemFs() {
	ctx := context.Background()
	p := authorized_clusters.K8sClusterConfig{
		ExtConfigStrID: "",
		ExtConfigID:    0,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			ClusterCfgStrID: "",
			CloudProvider:   "",
			Region:          "",
			Context:         "",
			Namespace:       "",
			Alias:           "",
			Env:             "",
		},
		IsActive: false,
		IsPublic: false,
	}
	_, err := CheckKubeConfig(ctx, s.Ou, p)
	s.Require().Nil(err)
}

func TestZeusAppCoreTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusAppCoreTestSuite))
}
