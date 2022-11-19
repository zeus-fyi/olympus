package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	test_base "github.com/zeus-fyi/olympus/test"
)

type ConfigMapTestSuite struct {
	suite.Suite
}

func (s *ConfigMapTestSuite) TestK8sConfigMapYamlReader() {
	cm := NewConfigMap()
	test_base.ForceDirToTestDirLocation()
	filepath := "./mocks/kubernetes_apps/demo/cm-demo.yaml"
	jsonBytes, err := ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &cm.K8sConfigMap)

	s.Require().Nil(err)
	s.Require().NotEmpty(cm.K8sConfigMap)
	s.Require().NotEmpty(cm.Metadata.Name)
	s.Require().NotEmpty(cm.Metadata.Labels)

	cm.ConvertK8sConfigMapToDB()
	s.Require().NotEmpty(cm.Data)
	s.Require().NotEmpty(cm.K8sConfigMap.Name)
	s.Require().NotEmpty(cm.K8sConfigMap.ObjectMeta.Name)
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

func TestConfigMapTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapTestSuite))
}
