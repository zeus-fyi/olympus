package create

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/sql_query"
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
	tmp := jen.Func().Params(jen.Id("s").Op("*").Id(structGen.Name)).Id(structGen.Name + "Insert")
	tmp.Add(tmpGenParams())
	tmp.Add(m.genFuncStructNameExamplesFieldCase())
	m.Add(tmp)
	return m.Save()
}

func (m *InsertModelTemplate) genFuncStructNameExamplesFieldCase() *jen.Statement {
	return jen.Block(m.genCompleteGenericExecSql()...)
}

func (m *InsertModelTemplate) genCompleteGenericExecSql() []jen.Code {
	return sql_query.GenPGGenericExec()
}

func tmpGenParams() *jen.Statement {
	return jen.Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("q").Id("sql_query_templates").Dot("QueryParams")).Params(jen.Id("error"))
}
