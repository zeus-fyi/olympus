package models

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/codegen/cookbook/core/template_test"
)

type StructNameExampleTestSuite struct {
	template_test.TemplateTestSuite
}

type Setup struct {
	template_test.TemplateTestSuite
}

func TestStructNameExampleTestSuite(t *testing.T) {
	suite.Run(t, new(StructNameExampleTestSuite))
}
