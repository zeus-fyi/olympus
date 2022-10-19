package funcs

import "github.com/zeus-fyi/jennifer/jen"

func (f *FuncGen) GetFuncHeader() *jen.Statement {
	headerStatement := jen.Id(f.Name)
	return headerStatement.Params(f.GetFieldStatement()...).Params(f.GetReturnFieldsStatement())
}
