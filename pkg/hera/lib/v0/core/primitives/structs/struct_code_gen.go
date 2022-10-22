package structs

import "github.com/zeus-fyi/jennifer/jen"

func (s *StructGen) GenerateStructJenCode() jen.Code {
	fs := make([]jen.Code, len(s.Fields))
	for i, f := range s.Fields {
		fs[i] = jen.Id(f.Name).Id(f.Type)
	}
	_struct := jen.Null().Type().Id(s.Name).Struct(fs...)
	return _struct
}

func (s *StructGen) GenerateStructJenStmt() *jen.Statement {
	fs := make([]jen.Code, len(s.Fields))
	for i, f := range s.Fields {
		fs[i] = jen.Id(f.Name).Id(f.Type)
	}
	_struct := jen.Null().Type().Id(s.Name).Struct(fs...)
	return _struct
}
