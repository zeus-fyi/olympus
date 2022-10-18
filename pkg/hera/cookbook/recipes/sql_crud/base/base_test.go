package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type ModelStructBaseGen struct {
	test.AutoGenBaseTestSuiteBase
}

var printOutLocation = "/Users/alex/Desktop/Zeus/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

func (s *ModelStructBaseGen) TestStructInsertFuncGen() {
	p := structs.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "model_template.go",
		Env:         "",
	}
	m := NewModelTemplate(p)
	//cg.Add(_struct.GenCreateStructWithExternalStructInheritance(wrapperStructName, extStructPath, extStructName))
	err := m.CreateTemplateFromStruct(StructMock())
	s.Require().Nil(err)
}

func TestModelStructBaseGen(t *testing.T) {
	suite.Run(t, new(ModelStructBaseGen))
}
