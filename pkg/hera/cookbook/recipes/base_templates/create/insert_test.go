package create

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
)

type StructInsertFuncGenRecipeTestSuite struct {
	cookbook.CookbookTestSuiteBase
}

func (s *StructInsertFuncGenRecipeTestSuite) TestStructInsertFuncGen() {
	fw := primitives.FileWrapper{}
	fw.PackageName = "autogen_structs"
	fw.FileName = "insert_model_template.go"
	err := GenStructPtrInsertFunc(fw)
	s.Require().Nil(err)
}

func TestStructInsertFuncGenRecipeTestSuite(t *testing.T) {
	suite.Run(t, new(StructInsertFuncGenRecipeTestSuite))
}
