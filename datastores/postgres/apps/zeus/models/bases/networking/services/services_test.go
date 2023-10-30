package services

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test/mocks"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (s *ServiceTestSuite) TestServiceParsing() {
	mocks.ChangeToMockDirectory()

	svc := NewService()
	b, err := ReadYamlConfig("./consensus_client/service.yaml")
	s.Require().Nil(err)
	s.Require().NotEmpty(b)

	err = json.Unmarshal(b, &svc.K8sService)
	s.Require().Nil(err)

	svc.ConvertK8sServiceToDB()
	s.Assert().NotEmpty(svc.ServicePorts)

	// override the clusterIP at the file location above then run this test
	// clusterIP: None
	//svc.ConvertServiceSpecConfigToDB()
	//s.Require().Equal("None", svc.ServiceSpec.ClusterIP.ChartSubcomponentValue)
}

func TestNetworkingTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func ReadYamlConfig(filepath string) ([]byte, error) {
	// Open YAML file
	jsonByteArray, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	jsonBytes, err := yaml.YAMLToJSON(jsonByteArray)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return jsonBytes, err
	}
	return jsonBytes, err
}
