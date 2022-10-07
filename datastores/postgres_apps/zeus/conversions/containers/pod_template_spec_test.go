package containers

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/test"
)

type TemplateSpecTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func TestTemplateSpecTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateSpecTestSuite))
}
