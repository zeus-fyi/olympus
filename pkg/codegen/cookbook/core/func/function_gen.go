package _func

import jen "github.com/dave/jennifer/jen"

func genDeclAt15() jen.Code {
	return jen.Null()
}
func genFunctemplateFunc() jen.Code {
	return jen.Func().Id("templateFunc").Params(jen.Id("ctx").Qual("context", "Context"),
		jen.Id("param").Id("string")).Params(jen.Id("error")).Block(jen.If(jen.Id("len").
		Call(jen.Id("param")).Op("<=").Lit(0)).
		Block(jen.Return().Qual("errors", "New").Call(jen.Lit("error message"))),
		jen.Return().Id("nil"))
}
func genFile() *jen.File {
	ret := jen.NewFile("core")
	ret.Add(genDeclAt15())
	ret.Add(genFunctemplateFunc())
	return ret
}
