package secrets

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type SecretsTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *SecretsTestSuite) TestK8sSecretYamlReaderAndK8sToDBCte() {
	secret := NewSecret()
	filepath := s.TestDirectory + "/mocks/test/secret.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &secret.K8sSecret)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret.K8sSecret)

	err = secret.ConvertK8sSecretToDB()
	s.Require().Nil(err)

	s.Require().NotEmpty(secret.Metadata)
	s.Require().NotEmpty(secret.StringData)
	s.Require().NotEmpty(secret.Type)

	c := charts.Chart{}
	c.ChartPackageID = 100
	cte := secret.GetSecretCTE(&c)
	s.Require().NotEmpty(cte)
}

func TestSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(SecretsTestSuite))
}
