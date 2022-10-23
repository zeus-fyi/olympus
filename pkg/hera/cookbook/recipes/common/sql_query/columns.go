package sql_query

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/funcs"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
)

func GeneratePgColumnsPtrFunc(structGen primitive.StructGen) jen.Code {
	fnName := "GetTableColumns"
	fn := funcs.NewFn(fnName)
	tableColVarName := "columnValues"
	columnValues := CreateStringSliceFieldAssignment(structGen)
	assignedColumnSlice := jen.Id(tableColVarName).Op("=").Index().Id("string").Values(columnValues...)
	fn.AddBodyStatement(assignedColumnSlice)
	returnField := fields.Field{
		Type: "[]string",
		Name: tableColVarName,
	}
	fn.AddReturnField(returnField)
	return fn.GenerateStructPtrFunc(structGen)
}

func GeneratePgTableNamePtrFunc(structGen primitive.StructGen) jen.Code {
	fnName := "GetTableName"
	fn := funcs.NewFn(fnName)
	tableName := "tableName"
	assignedColumnSlice := jen.Id(tableName).Op("=").Lit(structGen.TableExpressionName())
	fn.AddBodyStatement(assignedColumnSlice)
	returnField := fields.Field{
		Type: "string",
		Name: tableName,
	}
	fn.AddReturnField(returnField)
	return fn.GenerateStructPtrFunc(structGen)
}

func CreateStringSliceFieldAssignment(structGen primitive.StructGen) []jen.Code {
	var stmtChain []jen.Code
	cols := structGen.GetColumnFieldNames()
	for _, c := range cols {
		statement := &jen.Statement{}
		statement.Add(jen.Lit(c))
		stmtChain = append(stmtChain, statement)
	}
	return stmtChain
}
