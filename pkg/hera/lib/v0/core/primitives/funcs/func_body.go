package funcs

import "github.com/zeus-fyi/jennifer/jen"

func (f *FuncGen) AddBodyStatement(fncBodyElement *jen.Statement) {
	if len(f.Body) <= 0 {
		f.Body = []jen.Code{}
	}
	f.Body = append(f.Body, fncBodyElement)
	return
}

func (f *FuncGen) GetBodyAndReturn() *jen.Statement {
	f.AddBodyStatement(f.GetFuncReturnStatement())
	return jen.Block(f.Body...)
}
