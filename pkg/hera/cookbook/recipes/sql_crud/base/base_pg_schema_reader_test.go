package base

import "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"

func (s *ModelStructBaseGen) TestPGBaseSchemaReader() {
	p := structs.Path{
		PackageName: "autogen_bases",
		DirIn:       "",
		DirOut:      printOutLocation,
		Fn:          "model_template.go",
		Env:         "",
	}

	m := NewPGModelTemplate(p, nil, s.Tc.LocalDbPgconn)
	err := m.ReadPgTableDefinition()
	s.Require().Nil(err)
}
