package create_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	"k8s.io/apimachinery/pkg/util/rand"
)

type TopologyCreateClassRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *TopologyCreateClassRequestTestSuite) TestAddSkeletonBasesToCluster() {
	t.InitLocalConfigs()
	t.Eg.POST("/infra/class/skeleton/bases/create", CreateTopologySkeletonBasesActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	basesInsert := []string{"add-skeleton-base-" + rand.String(5), "add-skeleton-base-" + rand.String(5)}
	cc := zeus_req_types.TopologyCreateOrAddBasesToClassesRequest{
		ClassName:      "add-base-9l98z",
		ClassBaseNames: basesInsert,
	}

	_, err := t.ZeusClient.AddSkeletonBasesToClass(ctx, cc)
	t.Require().Nil(err)
}
func (t *TopologyCreateClassRequestTestSuite) TestAddBasesToCluster() {
	t.InitLocalConfigs()
	t.Eg.POST("/infra/class/bases/create", UpdateTopologyClassActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	basesInsert := []string{"add-base-" + rand.String(5), "add-base-" + rand.String(5)}
	cc := zeus_req_types.TopologyCreateOrAddBasesToClassesRequest{
		ClassName:      "unclassified-cluster",
		ClassBaseNames: basesInsert,
	}

	_, err := t.ZeusClient.AddBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *TopologyCreateClassRequestTestSuite) TestClassCreate() {
	t.InitLocalConfigs()
	t.Eg.POST("/infra/class/create", CreateTopologyClassActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	cc := zeus_req_types.TopologyCreateOrAddBasesToClassesRequest{
		ClassName: rand.String(10),
	}
	resp, err := t.ZeusClient.CreateClass(ctx, cc)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func TestTopologyCreateClassRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateClassRequestTestSuite))
}
