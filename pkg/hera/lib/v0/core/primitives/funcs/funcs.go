package funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

type FuncGen struct {
	Name         string
	Body         []jen.Code
	Fields       []fields.Field
	ReturnFields []fields.Field
}

func NewFn(name string) FuncGen {
	return FuncGen{
		Name:         name,
		Body:         []jen.Code{},
		Fields:       []fields.Field{},
		ReturnFields: []fields.Field{},
	}
}

func (f *FuncGen) GenerateFunc() *jen.Statement {
	fn := jen.Func()
	return fn.Add(f.GenerateFuncShell(fn))
}

func (f *FuncGen) GenerateFuncShell(prefix *jen.Statement) *jen.Statement {
	header := prefix.Add(f.GetFuncHeader())
	return header.Add(f.GetBodyAndReturn())
}
