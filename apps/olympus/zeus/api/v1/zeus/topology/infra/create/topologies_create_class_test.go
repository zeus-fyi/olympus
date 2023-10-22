package create_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateClassRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *TopologyCreateClassRequestTestSuite) TestRead() {
	apps.Pg.InitPG(context.Background(), t.Tc.ProdLocalDbPgconn)
	ctx := context.Background()
	cl, err := read_topology.SelectClusterTopology(ctx, 1679515557647002001, "avaxFujiNode", []string{"avax"})
	t.Require().Nil(err)
	t.Assert().NotEmpty(cl)
}

func (t *TopologyCreateClassRequestTestSuite) TestEndToEnd() {
	t.InitLocalConfigs()

	t.Eg.POST("/infra/create", CreateTopologyInfraActionRequestHandler)
	t.Eg.POST("/infra/class/skeleton/bases/create", CreateTopologySkeletonBasesActionRequestHandler)
	t.Eg.POST("/infra/class/bases/create", UpdateTopologyClassActionRequestHandler)
	t.Eg.POST("/infra/class/create", CreateTopologyClassActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()
	//
	//ctx := context.Background()
	//cc := zeus_req_types.TopologyCreateClusterClassRequest{
	//	ClusterClassName: rand.String(10),
	//}
	//fmt.Println(cc.ClusterClassName)
	//resp, err := t.ZeusClient.CreateClass(ctx, cc)
	//t.Require().Nil(err)
	//t.Assert().NotEmpty(resp)
	//
	//baseOne := "test-add-cluster-base-one" + rand.String(5)
	//baseTwo := "test-add-cluster-base-two" + rand.String(5)
	//basesInsert := []string{baseOne, baseTwo}
	//cb := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
	//	ClusterClassName:   cc.ClusterClassName,
	//	ComponentBaseNames: basesInsert,
	//}
	//
	//_, err = t.ZeusClient.AddComponentBasesToClass(ctx, cb)
	//t.Require().Nil(err)
	//
	//skBaseOne := "test-add-skeleton-base-one-" + rand.String(5)
	//skBaseTwo := "test-add-skeleton-base-two-" + rand.String(5)
	//skeletonBasesInsert := []string{skBaseOne}
	//cskb := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
	//	ClusterClassName:  cc.ClusterClassName,
	//	ComponentBaseName: baseOne,
	//	SkeletonBaseNames: skeletonBasesInsert,
	//}
	//_, err = t.ZeusClient.AddSkeletonBasesToClass(ctx, cskb)
	//skeletonBasesInsert2 := []string{skBaseTwo}
	//cskb2 := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
	//	ClusterClassName:  cc.ClusterClassName,
	//	ComponentBaseName: baseTwo,
	//	SkeletonBaseNames: skeletonBasesInsert2,
	//}
	//_, err = t.ZeusClient.AddSkeletonBasesToClass(ctx, cskb2)
	//cd := zeus_req_types.ClusterTopologyDeployRequest{
	//	ClusterClassName:    cc.ClusterClassName,
	//	SkeletonBaseOptions: []string{skBaseOne, skBaseTwo},
	//	CloudCtxNs:          beacon_cookbooks.BeaconCloudCtxNs,
	//}
	//
	//olympus_cookbooks.ChangeToCookbookDir()
	//c := beacon_cookbooks.ExecClientChart
	//p := beacon_cookbooks.BeaconExecClientChartPath
	//c.ClusterClassName = cc.ClusterClassName
	//c.ComponentBaseName = baseOne
	//c.SkeletonBaseName = skBaseOne
	//uploadResp, err := t.ZeusClient.UploadChart(ctx, p, c)
	//t.Require().Nil(err)
	//t.Assert().NotEmpty(uploadResp)
	//
	//c = beacon_cookbooks.ConsensusClientChart
	//p = beacon_cookbooks.BeaconConsensusClientChartPath
	//c.ClusterClassName = cc.ClusterClassName
	//c.ComponentBaseName = baseTwo
	//c.SkeletonBaseName = skBaseTwo
	//
	//uploadResp, err = t.ZeusClient.UploadChart(ctx, p, c)
	//t.Require().Nil(err)
	//t.Assert().NotEmpty(uploadResp)
	//
	//cl, err := read_topology.SelectClusterTopology(ctx, t.Tc.ProductionLocalTemporalOrgID, cd.ClusterClassName, cd.SkeletonBaseOptions)
	//t.Require().Nil(err)
	//t.Assert().NotEmpty(cl)
	//t.Assert().Len(cl.Topologies, 2)
	//t.Assert().Equal(cc.ClusterClassName, cl.ClusterClassName)
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
	//
	//basesInsert := []string{"add-skeleton-base-" + rand.String(5), "add-skeleton-base-" + rand.String(5)}
	//cc := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
	//	ClusterClassName:  "rqhppnzghs",
	//	ComponentBaseName: "add-base-9l98z",
	//	SkeletonBaseNames: basesInsert,
	//}
	//
	//fmt.Println(basesInsert)
	//_, err := t.ZeusClient.AddSkeletonBasesToClass(ctx, cc)
	//t.Require().Nil(err)
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

	//basesInsert := []string{"add-base-" + rand.String(5), "add-base-" + rand.String(5)}
	//cc := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
	//	ClusterClassName:   "rqhppnzghs",
	//	ComponentBaseNames: basesInsert,
	//}

	//_, err := t.ZeusClient.AddComponentBasesToClass(ctx, cc)
	//t.Require().Nil(err)
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
	//
	//cc := zeus_req_types.TopologyCreateClusterClassRequest{
	//	ClusterClassName: rand.String(10),
	//}
	//fmt.Println(cc.ClusterClassName)
	//resp, err := t.ZeusClient.CreateClass(ctx, cc)
	//t.Require().Nil(err)
	//t.Assert().NotEmpty(resp)
}

func TestTopologyCreateClassRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateClassRequestTestSuite))
}
