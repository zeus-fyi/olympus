package template_test

import jen "github.com/dave/jennifer/jen"

func genDeclAt24() jen.Code {
	return jen.Null()
}
func genDeclAt85() jen.Code {
	return jen.Null().Type().Id("FuncTestSuite").Struct(jen.Id("suite").Dot("Suite"))
}
func genFuncTesCodeGen() jen.Code {
	return jen.Func().Params(jen.Id("s").Op("*").Id("FuncTestSuite")).Id("TesCodeGen").Params().Block()
}
func genFuncTestFuncTestSuite() jen.Code {
	return jen.Func().Id("TestFuncTestSuite").Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(jen.Id("suite").Dot("Run").Call(jen.Id("t"), jen.Id("new").Call(jen.Id("FuncTestSuite"))))
}
func genFile() *jen.File {
	ret := jen.NewFile("template_test")
	ret.Add(genDeclAt24())
	ret.Add(genDeclAt85())
	ret.Add(genFuncTesCodeGen())
	ret.Add(genFuncTestFuncTestSuite())
	return ret
}
