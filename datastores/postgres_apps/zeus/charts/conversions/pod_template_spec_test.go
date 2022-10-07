package conversions

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TemplateSpecTestSuite struct {
	ChartPackagesTestSuite
}

func TestTemplateSpecTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateSpecTestSuite))
}
