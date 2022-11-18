package create

import (
	"testing"

	"github.com/stretchr/testify/suite"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type StructInsertFuncGenRecipeTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

var printOutLocation = "/Users/alex/Desktop/Zeus/olympus/datastores/postgres/apps/zeus/models/create/autogen"

func (s *StructInsertFuncGenRecipeTestSuite) TestStructInsertFuncGen() {
	p := filepaths.Path{
		PackageName: "autogen_structs",
		DirIn:       "",
		DirOut:      printOutLocation,
		FnIn:        "insert_model_template.go",
		Env:         "",
	}
	m := NewInsertModelTemplate(p)
	err := m.CreateTemplateFromStruct(InsertStructMock())
	s.Require().Nil(err)
}

func InsertStructMock() primitive.StructGen {
	structToMake := primitive.StructGen{
		Name:   "ChartPackageInsert",
		Fields: nil,
	}
	return structToMake
}

func TestStructInsertFuncGenRecipeTestSuite(t *testing.T) {
	suite.Run(t, new(StructInsertFuncGenRecipeTestSuite))
}
