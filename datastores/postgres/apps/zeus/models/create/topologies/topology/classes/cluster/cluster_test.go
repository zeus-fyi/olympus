package create_clusters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	create_systems "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/systems"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

var ctx = context.Background()

type CreateClustersTestSuite struct {
	test_suites.DatastoresTestSuite
}

const LocalEthereumBeaconClusterDefinitionID = 1670201797184939008
const UnclassifiedClusterDefinition = 0

func (s *CreateClustersTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

func (s *CreateClustersTestSuite) TestInsertClusterDefinition() {
	testDuplicate := "test-duplicate"
	c := NewClusterClassTopologyType(s.Tc.ProductionLocalTemporalOrgID, testDuplicate)
	err := create_systems.InsertSystem(ctx, &c.Systems)
	s.Require().Nil(err)
}

func (s *CreateClustersTestSuite) TestInsertCluster() {
	//name := "test-cluster"
	//err := delete_cluster.DeleteCluster(ctx, name)
	//
	//tx, err := apps.Pg.Begin(ctx)
	//s.Require().Nil(err)
	//
	//defer tx.Rollback(ctx)
	//sys := systems.Systems{TopologySystemComponents: autogen_bases.TopologySystemComponents{
	//	OrgID:                       s.Tc.ProductionLocalTemporalOrgID,
	//	TopologyClassTypeID:         class_types.ClusterClassTypeID,
	//	TopologySystemComponentName: name,
	//}}
	//
	//pcg := zeus_templates.ClusterPreviewWorkloads{
	//	ClusterName:    name,
	//	ComponentBases: make(map[string]map[string]topology_workloads.TopologyBaseInfraWorkload),
	//}
	//
	//ou := org_users.OrgUser{}
	//sbOne := make(map[string]topology_workloads.TopologyBaseInfraWorkload)
	//sbOne["sbTestBase1Test1"] = topology_workloads.TopologyBaseInfraWorkload{}
	//sbOne["sbTestBase1Test2"] = topology_workloads.TopologyBaseInfraWorkload{}
	//pcg.ComponentBases["baseTest1"] = sbOne
	//sbTwo := make(map[string]topology_workloads.TopologyBaseInfraWorkload)
	//sbTwo["sbTestBase2Test1"] = topology_workloads.TopologyBaseInfraWorkload{}
	//sbTwo["sbTestBase2Test2"] = topology_workloads.TopologyBaseInfraWorkload{}
	//pcg.ComponentBases["baseTest2"] = sbTwo
	//tx, err = InsertCluster(ctx, tx, &sys, pcg, ou)
	//s.Require().Nil(err)
	//
	//err = tx.Commit(ctx)
	//s.Require().Nil(err)
	//
	//tx, err = apps.Pg.Begin(ctx)
	//s.Require().Nil(err)
	//sbTwo["sbTestBase2Test3"] = topology_workloads.TopologyBaseInfraWorkload{}
	//tx, err = InsertCluster(ctx, tx, &sys, pcg, ou)
	//s.Require().Nil(err)
	//err = tx.Commit(ctx)
	//s.Require().Nil(err)
	//
	//tx, err = apps.Pg.Begin(ctx)
	//s.Require().Nil(err)
	//sbThree := make(map[string]topology_workloads.TopologyBaseInfraWorkload)
	//sbThree["sbTestBase3Test1"] = topology_workloads.TopologyBaseInfraWorkload{}
	//pcg.ComponentBases["baseTest3"] = sbThree
	//
	//tx, err = InsertCluster(ctx, tx, &sys, pcg, ou)
	//s.Require().Nil(err)
	//err = tx.Commit(ctx)

}

func TestCreateClustersTestSuite(t *testing.T) {
	suite.Run(t, new(CreateClustersTestSuite))
}
