package delete_cluster

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

var ctx = context.Background()

type DeleteClusterTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (t *DeleteClusterTestSuite) TestDeleteCluster() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	err := DeleteCluster(ctx, t.Tc.ProductionLocalTemporalOrgID, "pandora")
	t.Require().Nil(err)
	//cps := []string{"aws", "gcp", "do"}
	//networks := []string{"mainnet", "testnet", "devnet"}
	//
	//for _, net := range networks {
	//	for _, cp := range cps {
	//		name := fmt.Sprintf("sui-%s-%s", net, cp)
	//		oi := 1696626403975334000
	//		err := DeleteCluster(ctx, oi, name)
	//		t.Require().Nil(err)
	//	}
	//}
}

func TestDeleteClusterTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteClusterTestSuite))
}
