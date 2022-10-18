package funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

func (f *FuncGen) AddReturnField(field fields.Field) {
	f.ReturnFields = append(f.Fields, field)
}

func (f *FuncGen) AddReturnFields(field []fields.Field) {
	f.ReturnFields = append(f.Fields, field...)
}

// GetReturnFieldsStatement is for the variable name and data type
func (f *FuncGen) GetReturnFieldsStatement() *jen.Statement {
	statement := &jen.Statement{}
	for _, item := range f.ReturnFields {
		if len(item.Pkg) > 0 {
			statement.Add(jen.Id(item.Pkg).Dot(item.Type))
		} else {
			statement.Add(jen.Id(item.Name).Id(item.Type))
		}
	}
	return statement
}

// GetFuncReturnStatement is for an actual jen return with the variable names
func (f *FuncGen) GetFuncReturnStatement() *jen.Statement {
	returnStmt := jen.Return()
	for _, item := range f.ReturnFields {
		returnStmt.Add(jen.Id(item.Name))
	}
	return returnStmt
}
