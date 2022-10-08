package primitives

type StructGen struct {
	Name   string
	Fields []Field
}

func (s *StructGen) AddField(field Field) {
	s.Fields = append(s.Fields, field)
}

func (s *StructGen) AddFields(field []Field) {
	s.Fields = append(s.Fields, field...)
}
