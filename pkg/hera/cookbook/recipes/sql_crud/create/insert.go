package create

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type InsertModelTemplate struct {
	lib.CodeGen
}

func NewInsertModelTemplate(p structs.Path) InsertModelTemplate {
	m := InsertModelTemplate{lib.NewCodeGen(p)}
	return m
}

func (m *InsertModelTemplate) CreateTemplateFromStruct(structGen primitive.StructGen) error {
	m.Structs.AddStruct(structGen)
	m.Add(tmpGen(structGen.Name))
	return m.Save()
}

func tmpGen(structName string) jen.Code {
	tmp := jen.Func().Params(jen.Id("s").Op("*").Id(structName)).Id(structName + "Insert")
	tmp.Add(tmpGenParams())
	tmp.Add(genFuncStructNameExamplesFieldCase())
	return tmp
}

func tmpGenParams() *jen.Statement {
	return jen.Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("q").Id("sql_query_templates").Dot("QueryParams")).Params(jen.Id("error"))
}

func genLogHeader() *jen.Statement {
	return jen.Id("log").Dot("Debug").Call().Dot("Interface").Call(jen.Lit("InsertQuery:"),
		jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")))
}

func genFuncStructNameExamplesFieldCase() *jen.Statement {
	return jen.Block(genLogHeader(),
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
