package create_infra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulset"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/packages"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"k8s.io/apimachinery/pkg/util/rand"
)

type CreateInfraTestSuite struct {
	b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateInfraTestSuite) TestInsertInfraBase() {
	//nd := deployments.NewDeployment()
	nsvc := services.NewService()
	ing := ingresses.NewIngress()
	cm := configuration.NewConfigMap()
	sts := statefulset.NewStatefulSet()
	cw := chart_workload.ChartWorkload{
		//Deployment:  &nd,
		Service:     &nsvc,
		Ingress:     &ing,
		ConfigMap:   &cm,
		StatefulSet: &sts,
	}
	pkg := packages.Packages{
		Chart: charts.Chart{
			ChartPackages: autogen_bases.ChartPackages{
				ChartPackageID:   0,
				ChartName:        rand.String(10),
				ChartVersion:     rand.String(10),
				ChartDescription: sql.NullString{},
			},
		},
		ChartWorkload: cw,
	}

	//filepath := s.TestDirectory + "/apps/zeus/deployment.yaml"
	//jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	//s.Require().Nil(err)
	//err = json.Unmarshal(jsonBytes, &pkg.K8sDeployment)
	//s.Require().Nil(err)
	//err = pkg.ConvertDeploymentConfigToDB()
	//s.Require().Nil(err)

	filepath := s.TestDirectory + "/apps/zeus/service.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &pkg.K8sService)
	s.Require().Nil(err)
	pkg.ConvertK8sServiceToDB()
	s.Assert().NotEmpty(pkg.Service)

	filepath = s.TestDirectory + "/apps/zeus/cm-zeus.yaml"
	jsonBytes, err = s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &cm.K8sConfigMap)
	s.Require().Nil(err)
	cm.ConvertK8sConfigMapToDB()
	s.Assert().NotEmpty(cm.Data)
	s.Assert().NotEmpty(cm.Metadata.Name)

	filepath = s.TestDirectory + "/mocks/consensus_client/statefulset.yaml"
	jsonBytes, err = s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &sts.K8sStatefulSet)
	s.Require().Nil(err)
	err = sts.ConvertK8sStatefulSetToDB()
	s.Require().Nil(err)

	inf := NewCreateInfrastructure()
	inf.Packages = pkg

	ctx := context.Background()
	inf.Name = fmt.Sprintf("test_%d", s.Ts.UnixTimeStampNow())
	inf.OrgID, inf.UserID = s.b.NewTestOrgAndUser()
	err = inf.InsertInfraBase(ctx)
	s.Require().Nil(err)

	fmt.Println("ChartPackageID")
	fmt.Println(inf.ChartPackageID)
	fmt.Println("TopologyID")
	fmt.Println(inf.TopologyID)

}

func TestCreateInfraTestSuite(t *testing.T) {
	suite.Run(t, new(CreateInfraTestSuite))
}
