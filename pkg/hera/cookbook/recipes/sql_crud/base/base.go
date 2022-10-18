package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/funcs"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type ModelTemplate struct {
	lib.CodeGen
}

func NewModelTemplate(p structs.Path) ModelTemplate {
	m := ModelTemplate{lib.NewCodeGen(p)}
	return m
}

func (m *ModelTemplate) CreateTemplate() error {
	m.Structs.AddStruct(structMock())
	m.AddSlice(m.Structs.GenerateStructsJenCode(true))
	fn := "GetRowValues"
	m.Add(genFuncGetRowValues2(structMock(), fn))
	err := m.Save()
	return err
}

func createPgRowsVar() vars.VariableGen {
	v := vars.NewVarGen()
	v.InsertStruct(genRowValuesStructTemplate())
	return v
}

func genPgRowsFnFields() []fields.Field {
	return []fields.Field{{
		Name:  "queryName",
		Type:  "string",
		Value: "",
	}}
}

func genPgRowsFnReturnFields() []fields.Field {
	return []fields.Field{fields.Field{
		Pkg:   "apps",
		Name:  "pgValues",
		Type:  "RowValues",
		Value: "",
	}}
}

func genFuncGetRowValues2(structGen primitive.StructGen, fnName string) jen.Code {
	v := createPgRowsVar()
	declInitStruct := DeclarePgValuesStructVar(v, "init", v.GenStructInstructs["init"])
	statement := GenerateSwitchStatementForPgRows(v, structGen)

	fn := funcs.NewFn(fnName)
	fn.Fields = genPgRowsFnFields()
	fn.ReturnFields = genPgRowsFnReturnFields()

	fn.AddBodyStatement(declInitStruct)
	fn.AddBodyStatement(statement)

	return fn.GenerateStructPtrFunc(structGen)

}
