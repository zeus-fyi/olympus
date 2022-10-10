package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type ModelStructBaseGen struct {
	test.AutoGenBaseTestSuiteBase
}

var printOutLocation = "/Users/alex/Desktop/Zeus/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

func createTestCodeGenShell() lib.CodeGen {
	p := structs.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "model_template.go",
		Env:         "",
	}
	cg := lib.NewCodeGen(p)
	return cg
}

func (s *ModelStructBaseGen) TestStructInsertFuncGen() {
	cg := createTestCodeGenShell()
	//name := "ChartComponentKinds"
	//wrapperStructName := "Chart"
	//extStructPath := "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	//extStructName := "ChartPackages"

	//cg.Add(_struct.GenCreateStructWithExternalStructInheritance(wrapperStructName, extStructPath, extStructName))

	m := structMock()
	cg.Add(genHeader())
	cg.Add(genDeclAt85(m.Name))
	cg.Add(m.GenerateStructJenCode())
	cg.Add(m.GenerateSliceType())
	cg.Add(genFuncGetRowValues())
	cg.Add(m.GenerateStructJenCode())
	err := cg.Save()
	s.Require().Nil(err)
}

func TestModelStructBaseGen(t *testing.T) {
	suite.Run(t, new(ModelStructBaseGen))
}
