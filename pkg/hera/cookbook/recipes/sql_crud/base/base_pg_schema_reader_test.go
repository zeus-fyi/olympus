package base

import "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"

func (s *ModelStructBaseGenTestSuite) TestPGBaseSchemaReader() {
	p := filepaths.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocation,
		FnIn:        "model_template.go",
		Env:         "",
	}

	m := NewPGModelTemplate(p, nil, s.Tc.LocalDbPgconn)
	err := m.ReadPgTableDefinition()
	s.Require().Nil(err)
}
