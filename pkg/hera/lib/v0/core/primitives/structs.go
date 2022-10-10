package primitives

import "github.com/zeus-fyi/jennifer/jen"

type StructGen struct {
	Name   string
	Fields []Field
}

func (s *StructGen) AddField(field Field) {
	s.Fields = append(s.Fields, field)
}

func (s *StructGen) AddFields(field ...Field) {
	s.Fields = append(s.Fields, field...)
}

func (s *StructGen) GenerateStructJenCode() jen.Code {
	fields := make([]jen.Code, len(s.Fields))
	for i, f := range s.Fields {
		fields[i] = jen.Id(f.Name).Id(f.Type)
	}
	_struct := jen.Null().Type().Id(s.Name).Struct(fields...)
	return _struct
}

func (s *StructGen) GenerateSliceType() jen.Code {
	return jen.Null().Type().Id(s.Name + "s").Index().Id(s.Name)
}
