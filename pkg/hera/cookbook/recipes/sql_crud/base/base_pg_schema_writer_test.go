package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type ModelStructBaseGenWriterTestSuite struct {
	ModelStructBaseGenTestSuite
}

func (s *ModelStructBaseGenWriterTestSuite) TestPGBaseSchemaWriter() {

	filter := string_utils.FilterOpts{
		DoesNotStartWith: []string{"orgs", "user", "valid", "model"},
		StartsWith:       "",
		Contains:         "",
		DoesNotInclude:   nil,
	}
	p := structs.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "",
		Env:         "",
		FilterFiles: filter,
	}

	m := NewPGModelTemplate(p, nil, s.Tc.LocalDbPgconn)
	err := m.WritePgTableDefinition()
	s.Require().Nil(err)
}

func TestModelStructBaseGenWriterTestSuite(t *testing.T) {
	suite.Run(t, new(ModelStructBaseGenWriterTestSuite))
}
