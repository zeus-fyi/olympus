package funcs

import "github.com/zeus-fyi/jennifer/jen"

func (f *FuncGen) AddBodyStatement(fncBodyElement ...*jen.Statement) {
	if len(f.Body) <= 0 {
		f.Body = []jen.Code{}
	}
	for _, e := range fncBodyElement {
		f.Body = append(f.Body, e)
	}
	return
}

func (f *FuncGen) GetBodyAndReturn() *jen.Statement {
	f.AddBodyStatement(f.GetFuncReturnStatement())
	return jen.Block(f.Body...)
}
