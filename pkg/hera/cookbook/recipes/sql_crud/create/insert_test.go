package create

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/sql_crud/base"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type StructInsertFuncGenRecipeTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

var printOutLocation = "/Users/alex/Desktop/Zeus/olympus/datastores/postgres/apps/zeus/models/create/autogen"

func (s *StructInsertFuncGenRecipeTestSuite) TestStructInsertFuncGen() {
	p := structs.Path{
		PackageName: "autogen_structs",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "insert_model_template.go",
		Env:         "",
	}
	m := NewInsertModelTemplate(p)
	err := m.CreateTemplateFromStruct(base.StructMock())
	s.Require().Nil(err)
}

func TestStructInsertFuncGenRecipeTestSuite(t *testing.T) {
	suite.Run(t, new(StructInsertFuncGenRecipeTestSuite))
}
