package create

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	_struct "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/struct"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type StructInsertFuncGenRecipeTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

var printOutLocation = "/Users/alex/Desktop/Zeus/olympus/datastores/postgres/apps/zeus/structs/autogen_preview"

func createTestCodeGenShell() lib.CodeGen {
	p := structs.Path{
		PackageName: "autogen_structs",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "insert_model_template.go",
		Env:         "",
	}
	cg := lib.NewCodeGen(p)
	return cg
}
func (s *StructInsertFuncGenRecipeTestSuite) TestStructInsertFuncGen() {
	cg := createTestCodeGenShell()
	name := "ChartComponentKinds"
	wrapperStructName := "Chart"
	extStructPath := "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	extStructName := "ChartPackages"

	cg.Add(_struct.GenCreateStructWithExternalStructInheritance(wrapperStructName, extStructPath, extStructName))
	cg.Add(genDeclAt26(name))
	cg.Add(tmpGen(name))
	err := cg.Save()
	s.Require().Nil(err)
}

func TestStructInsertFuncGenRecipeTestSuite(t *testing.T) {
	suite.Run(t, new(StructInsertFuncGenRecipeTestSuite))
}
