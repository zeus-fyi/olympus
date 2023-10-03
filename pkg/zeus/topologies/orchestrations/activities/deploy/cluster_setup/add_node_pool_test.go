package deploy_topology_activities_create_setup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type DeployTestSuite struct {
	test_suites_base.TestSuite
}

func (s *DeployTestSuite) SetupTest() {
}

func (s *DeployTestSuite) TestMakeNodePoolReq() {

	act := CreateSetupTopologyActivities{}

	ctx := context.Background()
	params := base_deploy_params.ClusterSetupRequest{
		FreeTrial:  false,
		Ou:         org_users.OrgUser{},
		CloudCtxNs: zeus_common_types.CloudCtxNs{},
		Nodes: autogen_bases.Nodes{
			CloudProvider: "do",
			Region:        "nyc1",
			Slug:          "so1_5-16vcpu-128gb",
		},
		NodesQuantity: 1,
		Disks:         nil,
		Cluster:       zeus_templates.Cluster{},
		AppTaint:      false,
	}
	request, err := act.MakeNodePoolRequest(ctx, params)
	s.Require().Nil(err)
	s.Require().NotNil(request)
}

func TestDeployTestSuite(t *testing.T) {
	suite.Run(t, new(DeployTestSuite))
}

/*
{
  "FreeTrial": false,
  "Ou": {
    "orgID": 7138983863666904000,
    "userID": 7138958574876246000
  },
  "cloudProvider": "do",
  "region": "nyc1",
  "context": "do-nyc1-do-nyc1-zeus-demo",
  "namespace": "sui-testnet-do-d2ecb70d",
  "alias": "sui-testnet-do-d2ecb70d",
  "Nodes": {
    "memory": 0,
    "vcpus": 0,
    "disk": 0,
    "diskUnits": "",
    "diskType": "",
    "priceHourly": 0,
    "region": "nyc1",
    "cloudProvider": "do",
    "resourceID": 1680989079860309000,
    "description": "",
    "slug": "so1_5-16vcpu-128gb",
    "memoryUnits": "",
    "priceMonthly": 0,
    "gpus": 0,
    "gpuType": ""
  },
  "NodesQuantity": 1,
  "Disks": [
    {
      "resourceID": 1681408541855876000,
      "diskUnits": "2Ti",
      "priceMonthly": 0,
      "description": "",
      "type": "",
      "diskSize": 0,
      "priceHourly": 0,
      "region": "nyc1",
      "cloudProvider": "do"
    }
  ],
  "Cluster": {
    "clusterName": "sui-testnet-do",
    "componentBases": {
      "sui": {
        "sui": {
          "topologyID": "1696358364372655402",
          "addStatefulSet": true,
          "addDeployment": false,
          "addConfigMap": true,
          "addService": true,
          "addIngress": false,
          "addServiceMonitor": false,
          "configMap": null,
          "deployment": {
            "replicaCount": 0
          },
          "statefulSet": {
            "replicaCount": 0,
            "pvcTemplates": null
          },
          "containers": null,
          "resourceSums": {
            "replicas": "1",
            "memRequests": "63Gi",
            "memLimits": "63Gi",
            "cpuRequests": "7500m",
            "cpuLimits": "7500m",
            "diskRequests": "2Ti",
            "diskLimits": "0"
          }
        }
      },
      "suiIngress": {
        "suiIngress": {
          "topologyID": "1696304753265605122",
          "addStatefulSet": false,
          "addDeployment": false,
          "addConfigMap": false,
          "addService": false,
          "addIngress": true,
          "addServiceMonitor": false,
          "configMap": null,
          "deployment": {
            "replicaCount": 0
          },
          "statefulSet": {
            "replicaCount": 0,
            "pvcTemplates": null
          },
          "containers": null,
          "resourceSums": {
            "replicas": "",
            "memRequests": "",
            "memLimits": "",
            "cpuRequests": "",
            "cpuLimits": "",
            "diskRequests": "",
            "diskLimits": ""
          }
        }
      }
    },
    "ingressSettings": {
      "authServerURL": "",
      "host": ""
    },
    "ingressPaths": null
  },
  "AppTaint": true
}
*/
