package funcs

import "github.com/zeus-fyi/jennifer/jen"

func (f *FuncGen) GetFuncHeader() *jen.Statement {
	headerStatement := jen.Id(f.Name)
	headerStatement.Params(f.GetFieldStatement())
	return headerStatement
}
