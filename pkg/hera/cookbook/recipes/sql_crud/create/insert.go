package create

import (
	"github.com/zeus-fyi/jennifer/jen"
)

func tmpGen(structName string) jen.Code {
	tmp := jen.Func().Params(jen.Id("s").Op("*").Id(structName)).Id(structName + "Insert")
	tmp.Add(tmpGenParams())
	tmp.Add(genFuncStructNameExamplesFieldCase(structName))
	return tmp
}

func tmpGenParams() *jen.Statement {
	return jen.Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("q").Id("sql_query_templates").Dot("QueryParams")).Params(jen.Id("error"))
}
func genFuncStructNameExamplesFieldCase(structName string) *jen.Statement {
	return jen.Block(jen.Id("log").Dot("Debug").Call().Dot("Interface").Call(jen.Lit("InsertQuery:"),
		jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn"))),
		jen.List(jen.Id("r"), jen.Id("err")).Op(":=").Id("apps").
			Dot("Pg").Dot("Exec").Call(jen.Id("ctx"), jen.Id("q").Dot("SelectQuery").Call()),
		jen.If(jen.Id("returnErr").Op(":=").Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").
			Dot("LogHeader").Call(jen.Id("models").
			Dot("Sn"))), jen.Id("returnErr").Op("!=").Id("nil")).Block(jen.Return().Id("err")),
		jen.Id("rowsAffected").Op(":=").Id("r").Dot("RowsAffected").Call(),
		jen.Id("log").Dot("Debug").Call().Dot("Msgf").Call(jen.Lit("StructNameExamples: %s, Rows Affected: %d"),
			jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")), jen.Id("rowsAffected")),
		jen.Return().Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn"))))
}
func genDeclAt26(structName string) jen.Code {
	return jen.Null().Type().Id(structName).Struct(jen.Id("ChartComponentKindID").Id("int"), jen.Id("ChartComponentKindName").Id("string"), jen.Id("ChartComponentApiVersion").Id("string"))
}
