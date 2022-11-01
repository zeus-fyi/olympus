package create_infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type CreateInfraTestSuite struct {
	//b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateInfraTestSuite) TestInsertInfraBase() {
	//nd := deployments.NewDeployment()
	//nsvc := services.NewService()
	//ing := ingresses.NewIngress()
	//cm := configuration.NewConfigMap()
	//pkg := packages.Packages{
	//	Chart: charts.Chart{
	//		ChartPackages: autogen_bases.ChartPackages{
	//			ChartPackageID:   0,
	//			ChartName:        rand.String(10),
	//			ChartVersion:     rand.String(10),
	//			ChartDescription: sql.NullString{},
	//		},
	//	},
	//	Deployment: &nd,
	//	Service:    &nsvc,
	//	Ingress:    &ing,
	//	ConfigMap:  &cm,
	//}
	//
	//filepath := s.TestDirectory + "/apps/eth-indexer/deployment.yaml"
	//jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	//s.Require().Nil(err)
	//err = json.Unmarshal(jsonBytes, &pkg.K8sDeployment)
	//s.Require().Nil(err)
	//err = pkg.ConvertDeploymentConfigToDB()
	//s.Require().Nil(err)
	//
	//filepath = s.TestDirectory + "/apps/eth-indexer/service.yaml"
	//jsonBytes, err = s.Yr.ReadYamlConfig(filepath)
	//err = json.Unmarshal(jsonBytes, &pkg.K8sService)
	//s.Require().Nil(err)
	//pkg.ConvertK8sServiceToDB()
	//s.Assert().NotEmpty(pkg.Service)
	//
	//filepath = s.TestDirectory + "/apps/eth-indexer/ingress.yaml"
	//jsonBytes, err = s.Yr.ReadYamlConfig(filepath)
	//err = json.Unmarshal(jsonBytes, &pkg.K8sIngress)
	//s.Require().Nil(err)
	//err = pkg.ConvertK8sIngressToDB()
	//s.Require().Nil(err)
	//s.Assert().NotEmpty(pkg.Ingress)
	//
	//filepath = s.TestDirectory + "/apps/eth-indexer/cm-eth-indexer.yaml"
	//jsonBytes, err = s.Yr.ReadYamlConfig(filepath)
	//err = json.Unmarshal(jsonBytes, &cm.K8sConfigMap)
	//s.Require().Nil(err)
	//cm.ParseK8sConfigToDB()
	//s.Assert().NotEmpty(cm.Data)
	//s.Assert().NotEmpty(cm.Metadata.Name)
	//
	//inf := NewCreateInfrastructure()
	//inf.Packages = pkg
	//ctx := context.Background()
	//err = inf.InsertInfraBase(ctx)
	//s.Require().Nil(err)
}

func TestCreateInfraTestSuite(t *testing.T) {
	suite.Run(t, new(CreateInfraTestSuite))
}
