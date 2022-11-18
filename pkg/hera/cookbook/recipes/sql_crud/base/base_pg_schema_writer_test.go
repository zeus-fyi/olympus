package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type ModelStructBaseGenWriterTestSuite struct {
	ModelStructBaseGenTestSuite
}

var printOutLocationZeus = "/Users/alex/go/Olympus/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

func (s *ModelStructBaseGenWriterTestSuite) TestPGBaseSchemaWriter() {
	filter := string_utils.FilterOpts{
		DoesNotStartWithThese: []string{"orgs", "user", "valid", "model"},
		StartsWith:            "",
		Contains:              "",
		DoesNotInclude:        nil,
	}
	p := filepaths.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocationZeus,
		FnIn:        "",
		Env:         "",
	}

	m := NewPGModelTemplate(p, nil, s.Tc.LocalDbPgconn)
	m.Filter = &filter
	err := m.WritePgTableDefinition()
	s.Require().Nil(err)
}

var printOutLocationHestia = "/Users/alex/go/Olympus/olympus/datastores/postgres/apps/hestia/models/bases/autogen"

func (s *ModelStructBaseGenWriterTestSuite) TestHestiaBaseSchemaWriter() {
	filter := string_utils.FilterOpts{
		DoesNotStartWithThese: []string{},
		StartsWithThese:       []string{"org", "user"},
		Contains:              "",
		DoesNotInclude:        nil,
	}
	p := filepaths.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocationHestia,
		FnIn:        "",
		Env:         "",
	}

	m := NewPGModelTemplate(p, nil, s.Tc.LocalDbPgconn)
	m.Filter = &filter
	err := m.WritePgTableDefinition()
	s.Require().Nil(err)
}

func TestModelStructBaseGenWriterTestSuite(t *testing.T) {
	suite.Run(t, new(ModelStructBaseGenWriterTestSuite))
}
