package funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

//func tmpGenParams() *jen.Statement {
//	return jen.Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("q").Id("sql_query_templates").Dot("QueryParams")).Params(jen.Id("error"))
//}

func (f *FuncGen) GetFieldStatement() []jen.Code {
	var stmtChain []jen.Code
	for _, item := range f.Fields {
		statement := &jen.Statement{}
		if len(item.Pkg) > 0 {
			if len(item.Name) > 0 {
				statement.Add(jen.Id(item.Name).Id(item.Pkg).Dot(item.Type))
			} else {
				statement.Add(jen.Id(item.Pkg).Dot(item.Type))
			}
		} else {
			statement.Add(jen.Id(item.Name).Id(item.Type))
		}
		stmtChain = append(stmtChain, statement)
	}
	return stmtChain
}

func (f *FuncGen) AddField(field fields.Field) {
	f.Fields = append(f.Fields, field)
}

func (f *FuncGen) AddFields(field []fields.Field) {
	f.Fields = append(f.Fields, field...)
}
