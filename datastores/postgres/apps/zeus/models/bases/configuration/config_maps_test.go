package configuration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type ConfigMapTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ConfigMapTestSuite) TestK8sConfigMapYamlReader() {
	cm := NewConfigMap()
	filepath := s.TestDirectory + "/apps/eth-indexer/cm-eth-indexer.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &cm.K8sConfigMap)

	s.Require().Nil(err)
	s.Require().NotEmpty(cm.K8sConfigMap)
	s.Require().NotEmpty(cm.Metadata.Name)

	cm.ParseK8sConfigToDB()
	s.Require().NotEmpty(cm.Data)
	s.Require().NotEmpty(cm.K8sConfigMap.Name)
	s.Require().NotEmpty(cm.K8sConfigMap.ObjectMeta.Name)

}

func TestConfigMapTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapTestSuite))
}
