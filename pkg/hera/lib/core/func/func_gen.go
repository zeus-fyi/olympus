package _func

import (
	jen "github.com/dave/jennifer/jen"
	primitives2 "github.com/zeus-fyi/olympus/pkg/hera/lib/core/primitives"
)

func genFunctemplateFunc(fg primitives2.FuncGen) jen.Code {
	n := fg.Name
	ctxField := jen.Id("ctx").Qual("context", "Context")
	returnFuncParam := jen.Id("error")
	returnStatement := jen.Return().Id("nil")
	printLnStatement := jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world"))
	returnBlock := jen.Block(printLnStatement, returnStatement)
	fn := jen.Func().Id(n).Params(ctxField, fg.GetFieldStatement()).Params(returnFuncParam).Add(returnBlock)
	return fn
}

func genFile(fw primitives2.FileWrapper, funcGen primitives2.FuncGen) *jen.File {
	ret := jen.NewFile(fw.PackageName)
	ret.Add(genFunctemplateFunc(funcGen))
	return ret
}
