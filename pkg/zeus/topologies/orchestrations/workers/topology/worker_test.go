package topology_worker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type TopologyWorkerTestSuite struct {
	test_suites.TemporalTestSuite
}

func (s *TopologyWorkerTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

var ctx = context.Background()

func (s *TopologyWorkerTestSuite) TestExecuteDestroyClusterSetupWorkflowFreeTrial() {
	ta := s.Tc.DevTemporalAuth
	ns := "production-zeus.ngb72"
	hp := "production-zeus.ngb72.tmprl.cloud:7233"
	ta.Namespace = ns
	ta.HostPort = hp
	InitTopologyWorker(ta)
	cZ := Worker.ConnectTemporalClient()
	defer cZ.Close()
	Worker.Worker.RegisterWorker(cZ)
	err := Worker.Worker.Start()
	s.Require().Nil(err)
	params := base_deploy_params.DestroyClusterSetupRequest{
		ClusterSetupRequest: base_deploy_params.ClusterSetupRequest{
			FreeTrial: true,
			Ou:        org_users.NewOrgUserWithID(1696270811788972000, 1696270811788972000),
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "gcp",
				Region:        "us-central1",
				Context:       "gke_zeusfyi_us-central1-a_zeus-gcp-pilot-0",
				Namespace:     "f7194cf5-4dfa-408c-98d8-05615292a91a",
			},
			Nodes: hestia_autogen_bases.Nodes{
				Memory:        0,
				Vcpus:         0,
				Disk:          0,
				DiskUnits:     "",
				DiskType:      "",
				PriceHourly:   0,
				Region:        "us-central1",
				CloudProvider: "gcp",
				ResourceID:    1683165454211590000,
				Description:   "",
				Slug:          "e2-standard-8",
				MemoryUnits:   "",
				PriceMonthly:  0,
				Gpus:          0,
				GpuType:       "",
			},
			NodesQuantity: 1,
			Disks: hestia_autogen_bases.DisksSlice{
				{
					ResourceID:    1683165785839881000,
					Region:        "us-central1",
					CloudProvider: "gcp",
				},
			},
			Cluster: zeus_templates.Cluster{
				ClusterName: "microservice",
				ComponentBases: zeus_templates.ComponentBases{
					"microservice": zeus_templates.SkeletonBases{
						"api": zeus_templates.SkeletonBase{
							TopologyID:        "",
							AddStatefulSet:    false,
							AddDeployment:     true,
							AddConfigMap:      false,
							AddService:        true,
							AddIngress:        true,
							AddServiceMonitor: false,
							ConfigMap:         nil,
							Deployment:        zeus_templates.Deployment{},
							StatefulSet:       zeus_templates.StatefulSet{},
							Containers:        nil,
							ResourceSums:      autok8s_core.ResourceSums{},
						},
					},
				},
				IngressSettings: zeus_templates.Ingress{},
				IngressPaths:    nil,
			},
			AppTaint: true,
		},
	}
	err = Worker.ExecuteDestroyClusterSetupWorkflowFreeTrial(ctx, params)
	s.Require().Nil(err)
}

func TestTopologyWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyWorkerTestSuite))
}
