package servicemonitors

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type ServiceMonitorsTestSuite struct {
	test_suites_base.TestSuite
	TestDirectory string
}

func (s *ServiceMonitorsTestSuite) SetupTest() {
	s.TestDirectory = "./servicemonitor.yaml"
}
func (s *ServiceMonitorsTestSuite) TestStatefulSetK8sToDBConversion() {
	sm := NewServiceMonitor()
	filepath := s.TestDirectory
	jsonBytes, err := ReadYamlConfig(filepath)
	s.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &sm.K8sServiceMonitor)
	s.Require().Nil(err)
	s.Require().NotEmpty(sm.K8sServiceMonitor)
}

func TestServiceMonitorsTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceMonitorsTestSuite))
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
