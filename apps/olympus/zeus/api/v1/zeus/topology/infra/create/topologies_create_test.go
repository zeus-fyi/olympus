package create_infra

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestCreateTopology() {
	name := fmt.Sprintf("random_%d", t.Ts.UnixTimeStampNow())
	nd := deployments.NewDeployment()
	nsvc := services.NewService()
	ing := ingresses.NewIngress()
	cm := configuration.NewConfigMap()

	cw := chart_workload.ChartWorkload{
		Deployment: &nd,
		Service:    &nsvc,
		Ingress:    &ing,
		ConfigMap:  &cm,
	}
	pkg := packages.Packages{
		Chart:         charts.Chart{},
		ChartWorkload: cw,
	}
	t.c.TestDirectory = t.c.ForceDirToCallerLocation()

	filepath := t.c.TestDirectory + "/apps/eth-indexer/deployment.yaml"
	jsonBytes, err := t.c.Yr.ReadYamlConfig(filepath)
	t.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &pkg.K8sDeployment)
	t.Require().Nil(err)
	err = pkg.ConvertDeploymentConfigToDB()
	t.Require().Nil(err)

	filepath = t.c.TestDirectory + "/apps/eth-indexer/service.yaml"
	jsonBytes, err = t.c.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &pkg.K8sService)
	t.Require().Nil(err)
	pkg.ConvertK8sServiceToDB()
	t.Assert().NotEmpty(pkg.Service)

	filepath = t.c.TestDirectory + "/apps/eth-indexer/ingress.yaml"
	jsonBytes, err = t.c.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &pkg.K8sIngress)
	t.Require().Nil(err)
	err = pkg.ConvertK8sIngressToDB()
	t.Require().Nil(err)
	t.Assert().NotEmpty(pkg.Ingress)

	filepath = t.c.TestDirectory + "/apps/eth-indexer/cm-eth-indexer.yaml"
	jsonBytes, err = t.c.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &cm.K8sConfigMap)
	t.Require().Nil(err)

	oid, uid := t.h.NewTestOrgAndUser()
	orgUser := org_users.NewOrgUserWithID(oid, uid)

	c := charts.Chart{}
	c.ChartName = "test_api"
	c.ChartVersion = fmt.Sprintf("test_api_v%d", t.Ts.UnixTimeStampNow())
	nk := pkg.ChartWorkload.GetNativeK8s()

	tar := TopologyActionCreateRequest{
		TopologyActionRequest: base.CreateTopologyActionRequestWithOrgUser("create", orgUser),
		TopologyCreateRequest: TopologyCreateRequest{Name: name, Chart: c, NativeK8s: nk},
	}
	t.Endpoint = "/infra"
	t.AddEndpointHandler(tar.CreateTopology)
	tr := t.PostTopologyRequest(tar, 200)
	t.Require().NotEmpty(tr.Logs)

}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
