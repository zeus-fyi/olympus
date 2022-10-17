package funcs

import "github.com/zeus-fyi/jennifer/jen"

func (f *FuncGen) GetFuncBody() *jen.Statement {
	statement := &jen.Statement{}
	for _, item := range f.Body {
		statement.Add(item)
	}
	return statement
}

func (f *FuncGen) AddBodyStatement(fncBodyElement *jen.Statement) {
	if len(f.Body) <= 0 {
		f.Body = []*jen.Statement{}
	}
	f.Body = append(f.Body, fncBodyElement)
	return
}

func (f *FuncGen) GetFuncBodyAndReturn() *jen.Statement {
	statement := f.GetFuncBody()
	statement.Add(f.GetFuncReturnStatement())
	return statement
}
