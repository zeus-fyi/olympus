package base

import (
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/sql_query"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
)

func (m *ModelTemplate) CreateTemplateFromStruct(structGen primitive.StructGen) error {
	m.NewCodeGenFileShell()

	importNames := make(map[string]string)

	importNames["database/sql"] = "database/sql"
	// TODO not sure why, but for adding imports create a key and value for both pkg name and to itself
	importNames["github.com/zeus-fyi/olympus/datastores/postgres/apps"] = "github.com/zeus-fyi/olympus/datastores/postgres/apps"
	importNames["apps"] = "github.com/zeus-fyi/olympus/datastores/postgres/apps"

	m.JenFile.ImportNames(importNames)
	m.Structs.AddStruct(structGen)
	m.AddSlice(m.Structs.GenerateStructsJenCode(true))
	// these are template values
	v, structGen, bodyInitPgRowsStruct := GetPgRowsTemplateDeclarations(structGen)
	// each bodyPrefix variable is an independent body item in the function
	// you'll need to modify the generateSwitchStatementForPgRows fn to include more complex case conditions
	// it just uses a default of all rows for now
	bodySwitchStatement := generateSwitchStatementForPgRows(v, structGen)
	// you could add another body element here
	// fn template uses a default return type, the body is prefixed with body
	m.Add(GeneratePgRowsPtrFn(structGen, bodyInitPgRowsStruct, bodySwitchStatement))

	// adds table columns selector, todo refactor
	m.Add(sql_query.GeneratePgColumnsPtrFunc(structGen))

	// adds table name selector, todo refactor
	m.Add(sql_query.GeneratePgTableNamePtrFunc(structGen))

	err := m.Save()
	m.resetModelTemplate()
	return err
}

func (m *ModelTemplate) resetModelTemplate() {
	m.ResetInternalJenCaches()
	m.PrimitiveGenerator = primitives.PrimitiveGenerator{}
}
