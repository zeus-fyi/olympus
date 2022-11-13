package update_replicas

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type UpdateReplicasTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *UpdateReplicasTestSuite) TestUpdateStatefulSetReplicaCount() {
	s.InitLocalConfigs()
	r := ReplicaUpdate{}
	ctx := context.Background()

	rcCount := "2"
	topID, orgID, userID := 1668372506892811008, 7138983863666903883, 7138958574876245567
	trstart := read_topology.NewInfraTopologyReader()

	trstart.TopologyID = topID
	trstart.OrgID = orgID
	trstart.UserID = userID
	err := trstart.SelectTopology(ctx)
	s.Require().Nil(err)

	r.TopologyID = topID
	r.OrgID = orgID
	r.UserID = userID
	err = r.UpdateReplicaCountStatefulSet(ctx, rcCount)
	s.Require().Nil(err)

	tr := read_topology.NewInfraTopologyReader()

	tr.TopologyID = topID
	tr.OrgID = orgID
	tr.UserID = userID
	err = tr.SelectTopology(ctx)
	s.Require().Nil(err)

	s.Require().Equal(int32(2), *tr.K8sStatefulSet.Spec.Replicas)
	s.Assert().NotEqual(trstart.K8sStatefulSet.Labels["version"], tr.K8sStatefulSet.Labels["version"])
}

func (s *UpdateReplicasTestSuite) TestUpdateDeploymentReplicaCount() {
	s.InitLocalConfigs()
	r := ReplicaUpdate{}
	ctx := context.Background()

	rcCount := "4"
	topID, orgID, userID := 1668371829021945088, 7138983863666903883, 7138958574876245567
	trstart := read_topology.NewInfraTopologyReader()

	trstart.TopologyID = topID
	trstart.OrgID = orgID
	trstart.UserID = userID
	err := trstart.SelectTopology(ctx)
	s.Require().Nil(err)

	r.TopologyID = topID
	r.OrgID = orgID
	r.UserID = userID
	err = r.UpdateReplicaCountDeployment(ctx, rcCount)
	s.Require().Nil(err)

	tr := read_topology.NewInfraTopologyReader()

	tr.TopologyID = topID
	tr.OrgID = orgID
	tr.UserID = userID
	err = tr.SelectTopology(ctx)
	s.Require().Nil(err)

	s.Require().Equal(int32(4), *tr.K8sDeployment.Spec.Replicas)

	s.Assert().NotEqual(trstart.K8sDeployment.Labels["version"], tr.K8sDeployment.Labels["version"])
}

func TestUpdateReplicasTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateReplicasTestSuite))
}
