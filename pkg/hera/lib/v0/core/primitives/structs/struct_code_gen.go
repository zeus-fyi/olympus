package structs

import (
	"github.com/zeus-fyi/jennifer/jen"
)

func (s *StructGen) GenerateStructJenCode() jen.Code {
	fs := s.codeGenFieldsLoop()
	_struct := jen.Null().Type().Id(s.Name).Struct(fs...)
	return _struct
}

func (s *StructGen) GenerateStructJenStmt() *jen.Statement {
	fs := s.codeGenFieldsLoop()
	_struct := jen.Null().Type().Id(s.Name).Struct(fs...)
	return _struct
}

func (s *StructGen) codeGenFieldsLoop() []jen.Code {
	fs := make([]jen.Code, len(s.Fields))
	for i, f := range s.Fields {
		tags := f.GenerateTags()
		if len(tags) > 0 {
			fs[i] = jen.Id(f.Name).Id(f.Type).Tag(tags)
		} else {
			fs[i] = jen.Id(f.Name).Id(f.Type)
		}
	}
	return fs
}
