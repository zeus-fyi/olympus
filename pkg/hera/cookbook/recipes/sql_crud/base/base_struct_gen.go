package base

import primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"

func (m *ModelTemplate) CreateTemplateFromStruct(structGen primitive.StructGen) error {
	m.Structs.AddStruct(structGen)
	m.AddSlice(m.Structs.GenerateStructsJenCode(true))
	// these are template values
	v, structGen, bodyInitPgRowsStruct := GetPgRowsTemplateDeclarations()
	// each bodyPrefix variable is an independent body item in the function
	// you'll need to modify the generateSwitchStatementForPgRows fn to include more complex case conditions
	// it just uses a default of all rows for now
	bodySwitchStatement := generateSwitchStatementForPgRows(v, structGen)
	// you could add another body element here

	// fn template uses a default return type, the body is prefixed with body
	m.Add(GeneratePgRowsPtrFn(structGen, bodyInitPgRowsStruct, bodySwitchStatement))
	err := m.Save()
	return err
}
