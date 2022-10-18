package create

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/sql_query/common"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/sql_crud/base"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type InsertModelTemplate struct {
	base.ModelTemplate
}

func NewInsertModelTemplate(p structs.Path) InsertModelTemplate {
	sqlQueryType := "create"
	queryInfo := common.QueryMetadata{Type: sqlQueryType}
	m := InsertModelTemplate{base.NewModelTemplate(p, &queryInfo)}
	return m
}

func (m *InsertModelTemplate) CreateTemplateFromStruct(structGen primitive.StructGen) error {
	m.Structs.AddStruct(structGen)
	m.Add(m.tmpGen(structGen.Name))
	return m.Save()
}

func (m *InsertModelTemplate) tmpGen(structName string) jen.Code {
	tmp := jen.Func().Params(jen.Id("s").Op("*").Id(structName)).Id(structName + "Insert")
	tmp.Add(tmpGenParams())
	tmp.Add(m.genFuncStructNameExamplesFieldCase())
	return tmp
}

func (m *InsertModelTemplate) genFuncStructNameExamplesFieldCase() *jen.Statement {
	return jen.Block(genLogHeader(),
		jen.List(genSqlExec()),
		genSqlExecErrHandler(),
		genSqlExecRowsAffectedHandler(),
		genSqlExecRowsAffectedDebugLog(),
		genSqlInsertReturnWrappedErr())
}

func tmpGenParams() *jen.Statement {
	return jen.Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("q").Id("sql_query_templates").Dot("QueryParams")).Params(jen.Id("error"))
}

func genLogHeader() *jen.Statement {
	return jen.Id("log").Dot("Debug").Call().Dot("Interface").Call(jen.Lit("InsertQuery:"),
		jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")))
}

func genSqlExec() *jen.Statement {
	return jen.Id("err").Op(":=").Id("apps").Dot("Pg").Dot("Exec").Call(jen.Id("ctx"), jen.Id("q").Dot("SelectQuery").Call())
}

func genSqlExecErrHandler() *jen.Statement {
	return jen.If(jen.Id("returnErr").Op(":=").Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").
		Dot("LogHeader").Call(jen.Id("models").
		Dot("Sn"))), jen.Id("returnErr").Op("!=").Id("nil").Block(jen.Return().Id("err")))
}

func genSqlExecRowsAffectedHandler() *jen.Statement {
	return jen.Id("rowsAffected").Op(":=").Id("r").Dot("RowsAffected").Call()
}

func genSqlExecRowsAffectedDebugLog() *jen.Statement {
	return jen.Id("log").Dot("Debug").Call().Dot("Msgf").Call(jen.Lit("StructNameExamples: %s, Rows Affected: %d"),
		jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")), jen.Id("rowsAffected"))
}

func genSqlInsertReturnWrappedErr() *jen.Statement {
	return jen.Return().Id("misc").Dot("ReturnIfErr").Call(jen.Id("err"), jen.Id("q").Dot("LogHeader").Call(jen.Id("models").Dot("Sn")))
}
