package structs

import "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"

func (s *StructGen) AddField(field fields.Field) {
	s.Fields = append(s.Fields, field)
}

func (s *StructGen) AddFields(field ...fields.Field) {
	s.Fields = append(s.Fields, field...)
}

func (s *StructGen) GetColumnFieldNames() []string {
	var dbFieldColNames []string
	for _, f := range s.Fields {
		fn := f.DbFieldName()
		if len(fn) > 0 {
			dbFieldColNames = append(dbFieldColNames, f.DbFieldName())
		}
	}
	return dbFieldColNames
}
