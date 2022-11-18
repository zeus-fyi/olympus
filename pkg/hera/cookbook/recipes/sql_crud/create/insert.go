package create

import (
	"fmt"

	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/sql_query"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/sql_query/common"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/sql_crud/base"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/funcs"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type InsertModelTemplate struct {
	base.ModelTemplate
}

func NewInsertModelTemplate(p filepaths.Path) InsertModelTemplate {
	sqlQueryType := "create"
	queryInfo := common.QueryMetadata{Type: sqlQueryType}
	m := InsertModelTemplate{base.NewPGModelTemplate(p, &queryInfo, "")}
	return m
}

func (m *InsertModelTemplate) CreateTemplateFromStruct(structGen primitive.StructGen) error {
	m.Structs.AddStruct(structGen)
	m.Add(structGen.GenerateStructJenStmt())
	m.Add(m.GenerateModelPtrFn(structGen, sql_query.GenPGGenericExec()...))
	return m.Save()
}

// GenerateModelPtrFn generates boilerplate fn init
func (m *InsertModelTemplate) GenerateModelPtrFn(structGen primitive.StructGen, body ...*jen.Statement) jen.Code {
	sqlQueryName := fmt.Sprintf("%sInsert", structGen.Name)
	m.QueryMetadata.Name = sqlQueryName
	fn := funcs.NewFn(m.QueryMetadata.Name)
	fn.Fields = genInsertFnFields()
	fn.ReturnFields = genInsertFnReturnFields()
	fn.AddBodyStatement(body...)
	return fn.GenerateStructPtrFunc(structGen)
}
