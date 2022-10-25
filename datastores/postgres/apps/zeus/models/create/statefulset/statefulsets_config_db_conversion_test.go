package statefulset

import (
	"encoding/json"
	"testing"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	v1 "k8s.io/api/apps/v1"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

type ConvertStatefulSetPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ConvertStatefulSetPackagesTestSuite) TestConvertStatefulSet() {
	packageID := 0
	filepath := s.TestDirectory + "/mocks/test/statefulset.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var ss *v1.StatefulSet
	err = json.Unmarshal(jsonBytes, &ss)

	s.Require().Nil(err)
	s.Require().NotEmpty(ss)

	dbStatefulSetConfig, err := ConvertStatefulSetSpecConfigToDB(ss)
	s.Require().Nil(err)
	s.Require().NotEmpty(dbStatefulSetConfig)

	_ = dev_hacks.Use(packageID)
}

func TestConvertStatefulSetPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ConvertStatefulSetPackagesTestSuite))
}
