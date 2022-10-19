package funcs

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

func genFunctemplateFunc(fg FuncGen) jen.Code {
	n := fg.Name
	returnFuncParam := jen.Id("error")
	returnStatement := jen.Return().Id("nil")
	printLnStatement := jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world"))
	returnBlock := jen.Block(printLnStatement, returnStatement)
	fn := jen.Func().Id(n).Params(fg.GetFieldStatement()...).Params(returnFuncParam).Add(returnBlock)
	return fn
}

func genFile(fw fields.FileWrapper, funcGen FuncGen) *jen.File {
	ret := jen.NewFile(fw.PackageName)
	ret.Add(genFunctemplateFunc(funcGen))
	return ret
}
