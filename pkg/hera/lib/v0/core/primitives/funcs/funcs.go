package funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

type FuncGen struct {
	Name         string
	Body         []*jen.Statement
	Fields       []fields.Field
	ReturnFields []fields.Field
}

func (f *FuncGen) GenerateStructFunc(structName string) *jen.Statement {
	if len(structName) <= 0 {
		return jen.Nil()
	}
	shortHand := structName[0:1]
	fn := jen.Func().Params(jen.Id(shortHand).Op("*").Id(structName))
	return fn.Add(f.GenerateFuncShell())
}

func (f *FuncGen) GenerateFunc() *jen.Statement {
	fn := jen.Func()
	return fn.Add(f.GenerateFuncShell())
}

func (f *FuncGen) GenerateFuncShell() *jen.Statement {
	header := f.GetFuncHeader()
	bodyReturn := f.GetFuncBodyAndReturn()
	return header.Add(bodyReturn)
}
