package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type ModelStructBaseGenTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

var printOutLocation = "/Users/alex/go/Zeus/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

func (s *ModelStructBaseGenTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *ModelStructBaseGenTestSuite) TestBaseTemplateGeneration() {
	p := filepaths.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocation,
		FnIn:        "model_template.go",
		Env:         "",
	}

	m := NewPGModelTemplate(p, nil, s.Tc.LocalDbPgconn)
	//cg.Add(_struct.GenCreateStructWithExternalStructInheritance(wrapperStructName, extStructPath, extStructName))
	err := m.CreateTemplateFromStruct(StructMock())
	s.Require().Nil(err)
}

func TestModelStructBaseGen(t *testing.T) {
	suite.Run(t, new(ModelStructBaseGenTestSuite))
}
