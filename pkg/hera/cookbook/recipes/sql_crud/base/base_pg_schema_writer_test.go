package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type ModelStructBaseGenWriterTestSuite struct {
	ModelStructBaseGenTestSuite
}

func (s *ModelStructBaseGenWriterTestSuite) TestPGBaseSchemaWriter() {
	p := structs.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "model_template.go",
		Env:         "",
	}

	m := NewPGModelTemplate(p, nil, s.Tc.LocalDbPgconn)
	err := m.WritePgTableDefinition()
	s.Require().Nil(err)
}

func TestModelStructBaseGenWriterTestSuite(t *testing.T) {
	suite.Run(t, new(ModelStructBaseGenWriterTestSuite))
}
