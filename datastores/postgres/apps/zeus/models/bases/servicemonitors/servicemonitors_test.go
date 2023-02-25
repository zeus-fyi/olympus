package servicemonitors

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type ServiceMonitorsTestSuite struct {
	test_suites_base.TestSuite
	TestDirectory string
}

func (s *ServiceMonitorsTestSuite) SetupTest() {
	s.TestDirectory = "./servicemonitor.yaml"
}
func (s *ServiceMonitorsTestSuite) TestServiceMonitorK8sToDBConversion() {
	sm := NewServiceMonitor()
	filepath := s.TestDirectory
	jsonBytes, err := ReadYamlConfig(filepath)
	s.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &sm.K8sServiceMonitor)
	s.Require().Nil(err)
	s.Require().NotEmpty(sm.K8sServiceMonitor)

	err = sm.ConvertK8sServiceMonitorToDB()

	s.Require().Nil(err)
	s.Require().NotEmpty(sm.Metadata)
	c := charts.NewChart()
	ts := chronos.Chronos{}
	c.ChartPackageID = ts.UnixTimeStampNow()

	subCTEs := sm.GetServiceMonitorCTE(&c)
	s.Assert().NotEmpty(subCTEs)

	fmt.Println(subCTEs.GenerateChainedCTE())
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
