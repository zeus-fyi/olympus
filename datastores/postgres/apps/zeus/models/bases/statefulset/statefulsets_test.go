package statefulset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type StatefulSetTestSuite struct {
	base.TestSuite
	TestDirectory string
}

func (s *StatefulSetTestSuite) SetupTest() {
	s.TestDirectory = "./statefulset.yaml"
}
func (s *StatefulSetTestSuite) TestStatefulSetK8sToDBConversion() {
	sts := NewStatefulSet()
	filepath := s.TestDirectory
	jsonBytes, err := ReadYamlConfig(filepath)
	s.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &sts.K8sStatefulSet)
	s.Require().Nil(err)
	s.Require().NotEmpty(sts.K8sStatefulSet)

	err = sts.ConvertStatefulSetSpecConfigToDB()
	s.Require().Nil(err)
	s.Require().NotEmpty(sts.Spec)
	s.Require().NotEmpty(sts.Metadata)
	s.Require().NotEmpty(sts.Spec.Template)
	s.Require().NotEmpty(sts.Spec.Replicas)
	s.Require().NotEmpty(sts.Spec.Selector)
	s.Require().NotEmpty(sts.Spec.PodManagementPolicy)
	s.Require().NotEmpty(sts.Spec.VolumeClaimTemplates)
}

func TestStatefulSetTestSuite(t *testing.T) {
	suite.Run(t, new(StatefulSetTestSuite))
}

func ReadYamlConfig(filepath string) ([]byte, error) {
	// Open YAML file
	jsonByteArray, err := ioutil.ReadFile(filepath)
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
