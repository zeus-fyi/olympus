package base

import (
	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/sql_query/common"
	"github.com/zeus-fyi/olympus/pkg/hera/lib"
	primitive "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/vars"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type ModelTemplate struct {
	lib.CodeGen
	*common.QueryMetadata
}

// NewModelTemplate should TODO make the parent pkg name part more seamless
// Should use this to create all Model templates to avoid any nil panics on query metadata
func NewModelTemplate(p structs.Path, queryInfo *common.QueryMetadata) ModelTemplate {
	if queryInfo == nil {
		queryInfo = &common.QueryMetadata{Type: "base", ParentPkgName: "models"}
	}
	queryInfo.ParentPkgName = "models"
	m := ModelTemplate{lib.NewCodeGen(p), queryInfo}
	return m
}

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

// GenBaseImportHeaderLog to be used by structs which inherit this, shouldn't use for base logs without handling non-import pkgname prefix
func (m *ModelTemplate) GenBaseImportHeaderLog() *jen.Statement {
	modelVar := vars.NewVarGen()
	// modelVar.StringConstants["pkgName"] = "VarName" exported capital
	importedExternalVarName := "Sn"
	return common.GenPGDebugLogHeader(*m.QueryMetadata, modelVar.SetAndReturnImportedVarReference(m.ParentPkgName, importedExternalVarName))
}
