package structs

import (
	"strings"

	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

type StructGen struct {
	Name       string
	Fields     []fields.Field
	PluralDecl jen.Code
	//struct_sql_funcs.StructFuncGenSQL
}

func (s *StructGen) AddField(field fields.Field) {
	s.Fields = append(s.Fields, field)
}

func (s *StructGen) AddFields(field ...fields.Field) {
	s.Fields = append(s.Fields, field...)
}

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

func (s *StructGen) GenerateSliceType() jen.Code {
	s.PluralDecl = jen.Null().Type().Id(s.Name + "s").Index().Id(s.Name)
	return s.PluralDecl
}

func (s *StructGen) ShortHand() string {
	if len(s.Name) > 0 {
		return strings.ToLower(s.Name[0:1])
	}
	return ""
}
