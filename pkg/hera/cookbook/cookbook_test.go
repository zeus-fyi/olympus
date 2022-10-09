package cookbook

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	dirIn   = "models"
	dirOut  = "template_preview/models"
	fn      = "model_template.go"
	pkgName = "models"
	env     = "test"
)

type CookbookTestSuite struct {
	CookbookTestSuiteBase
}

func (s *CookbookTestSuite) TestAutoGenDir() {
	path := fileIO.NewFullPathDefinition(env, pkgName, dirIn, dirOut, fn)
	err := c.CreateTemplatesInPath(path)
	s.Require().Nil(err)
}

func (s *CookbookTestSuite) TestAutoGenDirTmpModels() {
	dirIn = "tmp_models"
	dirOut = "template_preview/tmp_models"

	path := fileIO.NewFullPathDefinition(env, pkgName, dirIn, dirOut, fn)
	err := c.CustomZeusParsing(path)
	s.Require().Nil(err)
}

func TestAutoCodeGenTestSuite(t *testing.T) {
	suite.Run(t, new(CookbookTestSuite))
}
