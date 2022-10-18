package funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

func (f *FuncGen) GetFieldStatement() *jen.Statement {
	statement := &jen.Statement{}
	for _, item := range f.Fields {
		if len(item.Pkg) > 0 {
			statement.Add(jen.Id(item.Pkg).Dot(item.Type))
		} else {
			statement.Add(jen.Id(item.Name).Id(item.Type))
		}
	}
	return statement
}

func (f *FuncGen) AddField(field fields.Field) {
	f.Fields = append(f.Fields, field)
}

func (f *FuncGen) AddFields(field []fields.Field) {
	f.Fields = append(f.Fields, field...)
}
