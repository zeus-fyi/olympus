package echo_api_app

import jen "github.com/dave/jennifer/jen"

func genDeclAt15() jen.Code {
	return jen.Null()
}
func genFuncmain() jen.Code {
	return jen.Func().Id("main").Params().Block(jen.Id("zerolog").Dot("TimeFieldFormat").Op("=").Id("zerolog").Dot("TimeFormatUnix"), jen.If(jen.Id("err").Op(":=").Id("server").Dot("ApiCmd").Dot("Execute").Call(), jen.Id("err").Op("!=").Id("nil")).Block(jen.Id("log").Dot("Err").Call(jen.Id("err"))))
}
func genFile() *jen.File {
	ret := jen.NewFile("main")
	ret.Add(genDeclAt15())
	ret.Add(genFuncmain())
	return ret
}
